package logPick

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"log-picker/kafka"
	"time"
)

// 通用日志采集配置
var config = tail.Config{
	ReOpen:    true,
	Follow:    true,
	Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	MustExist: false,
	Poll:      true,
}

type LogEtcd struct {
	File  string
	Topic string
}

type logAgentTask struct {
	tailObj    *tail.Tail
	logEtcd    LogEtcd
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (task *logAgentTask) init(logEtcd LogEtcd) {
	tails, err := tail.TailFile(logEtcd.File, config)
	if err != nil {
		fmt.Println("tail task failed, err:", err)
		return
	}
	task.tailObj = tails
	// 启动后台任务处理日志
	go task.run()
}

// 真实日志采集执行任务运行
func (task *logAgentTask) run() {
	for {
		select {
		case <-task.ctx.Done():
			fmt.Println("task done:", task.logEtcd.File, task.logEtcd.Topic)
			// 直接返回结束该任务
			return
		case line, ok := <-task.tailObj.Lines:
			if !ok {
				fmt.Printf("tail task close reopen, filename:%s\n", task.tailObj.Filename)
				time.Sleep(time.Second)
				continue
			}
			// 发送到kafka管道
			kafka.SendToChan(task.logEtcd.Topic, line.Text)
		}
	}
}

// runAgentTask 新建日志采集任务
func runAgentTask(logEtcd LogEtcd) *logAgentTask {
	ctx, cancel := context.WithCancel(context.Background())
	task := &logAgentTask{
		logEtcd:    logEtcd,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	task.init(logEtcd)
	fmt.Println("new task:", logEtcd.File, logEtcd.Topic)
	// 新增一个消费任务
	runTransferTask(logEtcd, ctx, cancel)
	return task
}
