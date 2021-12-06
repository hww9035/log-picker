package logPick

import (
	"time"
)

type taskManager struct {
	logEtcdList      []*LogEtcd
	agentTasks       map[string]*logAgentTask
	newAgentTaskChan chan []*LogEtcd
}

// 全局日志采集管理器
var tskM *taskManager

// Init 初始化日志管理器
func Init(list []*LogEtcd) {
	tskM = &taskManager{
		logEtcdList:      list,
		agentTasks:       make(map[string]*logAgentTask),
		newAgentTaskChan: make(chan []*LogEtcd),
	}
	for _, v := range list {
		task := runAgentTask(*v)
		tskM.agentTasks[v.File] = task
	}
	go checkTaskStatus()
}

// 监听日志采集配置变化
func checkTaskStatus() {
	for {
		select {
		case nt := <-tskM.newAgentTaskChan:
			tskM.logEtcdList = nt
			for _, log := range nt {
				t, ok := tskM.agentTasks[log.File]
				if !ok {
					tskM.agentTasks[log.File] = runAgentTask(*log)
				}
				if ok && t.logEtcd.Topic != log.Topic {
					t.cancelFunc()
					tskM.agentTasks[log.File] = runAgentTask(*log)
				}
			}
			isDel := true
			for _, ot := range tskM.agentTasks {
				for _, newt := range nt {
					if newt.File == ot.logEtcd.File && newt.Topic == ot.logEtcd.Topic {
						isDel = false
						continue
					}
				}
				if isDel {
					// 通知子任务结束
					ot.cancelFunc()
				}
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// AddNewTaskChan 暴露给外部调用，通知管理器配置信息有变化
func AddNewTaskChan(newTask []*LogEtcd) {
	tskM.newAgentTaskChan <- newTask
}
