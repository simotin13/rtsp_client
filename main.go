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

	conn, err := net.Dial("tcp", os.Args[0])
	if err != nil {
		os.Exit(-1)
	}
	defer conn.Close()

	url := os.Args[0]

	// OPTIONS
	req := fmt.Sprintf("OPTIONS %s RTSP/1.0\r\n", url)
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
	if len(ary) < 4 {
		return false
	}
	if ary[0] != "Replay:" {
		return false
	}

	if ary[1] != "RTSP/1.0" {
		return false
	}

	code, err := strconv.Atoi(ary[2])
	if err != nil {
		return false
	}
	if code != 200 {
		return false
	}
	if ary[3] != "OK" {
		return false
	}

	// success only 200 OK
	return true
}
