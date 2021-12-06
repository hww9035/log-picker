package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

var (
	client    *elastic.Client
	msgEsChan chan *msgEs
)

type msgEs struct {
	topic string
	msg   string
}

// Init 初始化es
func Init(address string) (err error) {
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err == nil {
		// 初始化消息管道
		msgEsChan = make(chan *msgEs)
		// 后台异步处理发送消息到es
		go sendToEs()
	}
	return err
}

// 从管道获取消息发送到es
func sendToEs() error {
	for {
		select {
		case msg := <-msgEsChan:
			_, err := client.Index().Index(msg.topic).BodyJson(*msg).Do(context.Background())
			if err != nil {
				fmt.Println("send to es fail:", err)
			} else {
				fmt.Println("send to es ok:", msg)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// SendEsChan 发送消息到管道
func SendEsChan(topic string, msg string) {
	msgEsChan <- &msgEs{
		topic: topic,
		msg:   msg,
	}
}
