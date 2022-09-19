package network

import (
	"bufio"
	"fmt"
	"net"
)

var addressListen = "127.0.0.1:20000"

func process(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		// 读取数据并写入buf中，返回读取的字节数
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
		// 发送数据
		conn.Write([]byte(recvStr))
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
