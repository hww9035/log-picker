package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log-picker/mq/nsq"
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
	_ = nsq.InitProducer("127.0.0.1:14150")
	go func() {
		for {
			t := time.Now().Unix()
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-1", t))
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-2", t))
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-3", t))
			time.Sleep(time.Second * 3)
		}
	}()
	nsq.TestConsumer("127.0.0.1:4161", "top1", "chan1")
}
