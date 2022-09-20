package network

import (
	"bufio"
	"fmt"
	"net"
)

var addressListen = "127.0.0.1:20000"

func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		receiveStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", receiveStr)
		_, _ = conn.Write([]byte(receiveStr))
	}
}

func testServer() {
	listen, err := net.Listen("tcp", addressListen)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
