package logPick

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log-picker/es"
	"log-picker/kafka"
)

type logTransferTask struct {
	logEtcd    LogEtcd
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (task *logTransferTask) run() {
	addEs(task)
}

func addEs(task *logTransferTask) {
	consumer := kafka.GetConsumer()
	// 根据topic取到所有的分区
	partitionList, err := consumer.Partitions(task.logEtcd.Topic)
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	for partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(task.logEtcd.Topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		go sendEs(task, pc)
	}
}

func sendEs(task *logTransferTask, pc sarama.PartitionConsumer) {
	defer pc.AsyncClose()
	// 消费者管道获取消息

	//for msg := range pc.Messages() {
	//	es.SendEsChan(task.logEtcd.Topic, string(msg.Value))
	//}

	for {
		select {
		case <-task.ctx.Done():
			return
		case msg := <-pc.Messages():
			es.SendEsChan(task.logEtcd.Topic, string(msg.Value))
		}
	}
}

func runTransferTask(cf LogEtcd, ctx context.Context, cancel context.CancelFunc) {
	task := &logTransferTask{
		logEtcd:    cf,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	fmt.Println("new transfer task:", cf.File, cf.Topic)
	task.run()
}
