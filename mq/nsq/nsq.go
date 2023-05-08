package nsq

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	nsqProducer *nsq.Producer
)

func InitProducer(address string) (err error) {
	config := nsq.NewConfig()
	nsqProducer, err = nsq.NewProducer(address, config)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func PubMsg(topicName string, msg string) (err error) {
	messageBody := []byte(msg)
	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	err = nsqProducer.Publish(topicName, messageBody)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// nsqProducer.Stop()
	return err
}

type myMessageHandler struct{}

func (h *myMessageHandler) HandleMessage(msg *nsq.Message) error {
	if len(msg.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	fmt.Printf("recv from %v, msg: %v\n", msg.NSQDAddress, string(msg.Body))

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	// return errors.New("fail")
	return nil
}

func TestConsumer(lookupdAddress, topic, channel string) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	nsqConsumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal("create consumer failed, err:", err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	nsqConsumer.AddHandler(&myMessageHandler{})

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = nsqConsumer.ConnectToNSQLookupd(lookupdAddress)
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("consumer stop by signal:", sig)

	// Gracefully stop the consumer.
	nsqConsumer.Stop()
}
