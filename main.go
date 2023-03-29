package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("start client")
	if len(os.Args) < 1 {
		fmt.Println("Usage rtsp_client rtsp://<ipaddr>:<port>/<path>")
		os.Exit(1)
	}
	url := os.Args[1]

	// rtsp://user:passwd@0.0.0.0:port/path
	// rtsp://0.0.0.0
	idx := strings.Index(url, "rtsp://")
	if idx < 0 {
		// invalid url
		fmt.Sprintf("invalud url [%s]", url)
		os.Exit(1)
	}

	tcpAddr := ""
	user := ""
	passwd := ""
	tmp := url[7:]
	idx = strings.Index(tmp, "@")
	if idx < 0 {
		// ユーザ名:パスワード なし
	} else {
		// ユーザ名:パスワード あり
		userPass := tmp[:idx]
		ary := strings.Split(userPass, ":")
		if len(ary) < 2 {
			fmt.Printf("invalud url [%s]", url)
			os.Exit(1)
		}
		user = ary[0]
		passwd = ary[1]
		fmt.Printf("%s,%s", user, passwd)
		tmp = tmp[(idx + 1):]
		url = "rtsp://" + tmp
	}

	// ホスト名・IPアドレス部分を切り取り
	idx = strings.LastIndex(tmp, "/")
	if 0 < idx {
		tcpAddr = tmp[0:idx]
	}

	if !strings.Contains(tcpAddr, ":") {
		// RTSPのデフォルトポート番号を指定
		tcpAddr += ":554"
	}
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer conn.Close()

	// OPTIONS
	// OPTIONS rtsp://192.168.1.222:554/ipcam_h264.sdp RTSP/1.0
	//CSeq: 2
	// User-Agent: LibVLC/3.0.18 (LIVE555 Streaming Media v2016.11.28)
	//OPTIONS rtsp://192.168.1.222:554/ipcam_h264.sdp RTSP/1.0
	// CSeq: 2
	//User-Agent: LibVLC/3.0.18 (LIVE555 Streaming Media v2016.11.28)

	seqNo := 2
	req := fmt.Sprintf("OPTIONS %s RTSP/1.0\r\n", url)
	req += fmt.Sprintf("CSeq:%d\r\n", seqNo)
	seqNo++
	req += "User-Agent: MyRTSP client\r\n"
	req += "\r\n"
	resp, err := sendRecv(conn, req)
	if err != nil {
		os.Exit(-2)
	}

	if !isOk(resp) {
		msg := fmt.Sprintf("OPTIONS failed, resp:[%s]", resp)
		fmt.Println(msg)
		os.Exit(-2)
	}

	// DESCRIBE
	req = fmt.Sprintf("DESCRIBE %s RTSP/1.0\r\n", url)
	resp, err = sendRecv(conn, req)
	if err != nil {
		os.Exit(-2)
	}

	if !isOk(resp) {
		msg := fmt.Sprintf("DESCRIBE failed, resp:[%s]", resp)
		fmt.Println(msg)
		os.Exit(-2)
	}

	// CSeqの解析
	// レスポンスボディの解析

	// SETUP
	req = fmt.Sprintf("SETUP %s/trackID=0 RTSP/1.0\r\n", url)
	resp, err = sendRecv(conn, req)
	if err != nil {
		os.Exit(-2)
	}

	if !isOk(resp) {
		msg := fmt.Sprintf("DESCRIBE failed, resp:[%s]", resp)
		fmt.Println(msg)
		os.Exit(-2)
	}
}

func sendRecv(conn net.Conn, request string) (string, error) {
	resp := ""
	b := []byte(request)
	_, err := conn.Write(b)
	if err != nil {
		return resp, err
	}

	resp, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return resp, err
	}

	resp = strings.TrimRight(resp, "\r\n")
	return resp, nil
}

func isOk(resp string) bool {
	ary := strings.Split(resp, " ")
	if len(ary) < 3 {
		return false
	}

	if ary[0] != "RTSP/1.0" {
		return false
	}

	code, err := strconv.Atoi(ary[1])
	if err != nil {
		return false
	}
	if code != 200 {
		return false
	}
	if ary[2] != "OK" {
		return false
	}

	// success only 200 OK
	return true
}
