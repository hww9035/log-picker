package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type msgKafkaData struct {
	topic string
	msg   string
}

var (
	clientProducer sarama.SyncProducer
	clientConsumer sarama.Consumer
	msgKafkaChan   chan *msgKafkaData
)

// InitProducer 初始化生产者
func InitProducer(addr string, chanSize int) (err error) {
	// 生产者配置
	config := sarama.NewConfig()
	// 发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true
	clientProducer, err = sarama.NewSyncProducer([]string{addr}, config)
	if err == nil {
		// 初始化成功则构建消息管道
		msgKafkaChan = make(chan *msgKafkaData, chanSize)
		go sendToKafka()
	}
	return err
}

// InitConsumer 初始化消费者
func InitConsumer(addr string) (err error) {
	clientConsumer, err = sarama.NewConsumer([]string{addr}, nil)
	return err
}

func GetProducer() sarama.SyncProducer {
	return clientProducer
}

func GetConsumer() sarama.Consumer {
	return clientConsumer
}

func sendToKafka() {
	for {
		select {
		case data := <-msgKafkaChan:
			// 构造一个消息并发送
			msg := &sarama.ProducerMessage{
				Topic: data.topic,
				Value: sarama.StringEncoder(data.msg),
			}
			_, _, err := clientProducer.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg to kafka failed, err:", err)
			} else {
				fmt.Println("send msg to kafka ok:", msg.Topic, data.msg)
			}
			// 客户端查看具体topic内容：
			// kafka/bin/kafka-console-consumer.sh --topic topicName --from-beginning --bootstrap-server localhost:9092
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func SendToChan(topic string, msg string) {
	if strings.TrimSpace(msg) == "" {
		return
	}
	msgKafkaChan <- &msgKafkaData{
		topic: topic,
		msg:   msg,
	}
}
