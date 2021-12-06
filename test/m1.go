package test

import (
	"encoding/json"
	"fmt"
	"log-picker/etcd"
	"log-picker/logPick"
)

func addLogAgentEtcd() {
	err := etcd.Init("127.0.0.1:2379", 5)
	if err != nil {
		fmt.Println("etcd err:", err)
		return
	}
	msgMap := make([]logPick.LogEtcd, 0)
	msgMap = append(msgMap, logPick.LogEtcd{
		File:  "/Users/huangweiwei/mysql.log",
		Topic: "mysql-log",
	})
	msgMap = append(msgMap, logPick.LogEtcd{
		File:  "/Users/huangweiwei/redis.log",
		Topic: "redis-log",
	})
	msg, err := json.Marshal(msgMap)
	if err != nil {
		fmt.Println("Marshal fail:", err)
		return
	}
	msgStr := string(msg)
	fmt.Println(msgStr)
	ok := etcd.Put("log_agent", msgStr)
	if !ok {
		fmt.Println("put fail")
		return
	}
	etcd.CloseEtcd()
	fmt.Println("put ok")
}

func getLogAgentEtcd() {
	err := etcd.Init("127.0.0.1:2379", 5)
	if err != nil {
		fmt.Println("etcd err:", err)
		return
	}
	d, err := etcd.Get("log_agent")
	if err != nil {
		fmt.Println("get fail")
		return
	}
	etcd.CloseEtcd()
	fmt.Println(d)
}
