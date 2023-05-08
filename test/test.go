package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"

	"log-picker/queue"
)

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
	_ = buf.Flush()
	// 方式二：copy
	// n, err := io.Copy(file, res.Body)

	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("down size:", n/1024/1024, " M")
}

func TestNsq() {
	_ = queue.InitProducer("127.0.0.1:4150")
	_ = queue.PubMsg("top1", "hello1")
	_ = queue.PubMsg("top1", "hello2")
	_ = queue.PubMsg("top1", "hello3")
	time.Sleep(time.Second * 3)

	queue.TestConsumer("127.0.0.1:4161", "top1", "chan1")
}
