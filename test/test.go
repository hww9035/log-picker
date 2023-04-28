package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Btest() {
	fmt.Println("Btest")
}

func Down() {
	file, err := os.Create("/Users/huangweiwei/Downloads/go.tar.gz")
	if err != nil {
		log.Fatal(err)
		return
	}
	res, _ := http.Get("https://go.dev/dl/go1.20.3.darwin-amd64.tar.gz")

	// 方式一：buf
	buf := bufio.NewWriter(file)
	n, err := buf.ReadFrom(res.Body)
	// 方式二：copy
	//n, err := io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	_ = buf.Flush()

	log.Println("down size:", n/1024/1024)
}
