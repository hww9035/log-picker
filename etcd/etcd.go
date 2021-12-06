package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log-picker/logPick"
	"time"
)

var (
	client *clientv3.Client
)

// Init 初始化etcd连接
func Init(addr string, timeout int) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	return err
}

func CloseEtcd() {
	client.Close()
}

func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return "", err
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			return string(ev.Value), nil
		}
	}
	return "", err
}

func GetConf(key string) ([]*logPick.LogEtcd, error) {
	cfg, err := Get(key)
	if err != nil {
		fmt.Println("get conf fail:", err)
		return nil, err
	}
	var logEtcdSlice = make([]*logPick.LogEtcd, 0)
	json.Unmarshal([]byte(cfg), &logEtcdSlice)
	return logEtcdSlice, nil
}

func Put(key, val string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := client.Put(ctx, key, val)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return false
	}
	return true
}

func Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := client.Delete(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("del from etcd failed, err:%v\n", err)
	}
	return err
}

func Watch(key string) {
	// WatchChan <-chan WatchResponse
	rch := client.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			var logEtcdSlice = make([]*logPick.LogEtcd, 0)
			json.Unmarshal(ev.Kv.Value, &logEtcdSlice)
			if ev.Type == 0 {
				// PUT
				fmt.Println("etcd config put:", string(ev.Kv.Value))
			}
			if ev.Type == 1 {
				// DELETE
				fmt.Println("etcd config delete:", string(ev.Kv.Value))
			}
			// 通知日志采集管理器配置出现变化
			logPick.AddNewTaskChan(logEtcdSlice)
		}
	}
}
