package main

import (
    "log"

    "gopkg.in/ini.v1"
    "log-picker/conf"
    "log-picker/es"
    "log-picker/etcd"
    "log-picker/kafka"
    "log-picker/logPick"
)

func main() {
    // 加载基础配置
    var cfg conf.Config
    err := ini.MapTo(&cfg, "conf/conf.ini")
    if err != nil {
        log.Fatal(err)
    }

    // mysql初始化
    //err = db.Init(cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.User, cfg.Mysql.Pwd, cfg.Mysql.DbName)
    //if err != nil {
    //    log.Fatal(err)
    //}

    // redis初始化
    //err = redis.Init(cfg.Redis.Address)
    //if err != nil {
    //    log.Fatal(err)
    //}

    // etcd初始化，监听配置变化big传递给后续日志收集服务
    err = etcd.Init(cfg.Endpoints, cfg.Etcd.DialTimeout)
    if err != nil {
        log.Fatal(err)
    }

    // kafka初始化，其中后台等待发送处理消息
    err = kafka.InitProducer(cfg.Kafka.Address, cfg.Kafka.ChanSize)
    if err != nil {
        log.Fatal(err)
    }
    err = kafka.InitConsumer(cfg.Kafka.Address)
    if err != nil {
        log.Fatal(err)
    }

    //es初始化
    err = es.Init(cfg.Es.Address)
    if err != nil {
        log.Fatal(err)
    }

    // logAgent初始化
    logEtcdList, err := etcd.GetConf(cfg.Etcd.Key)
    if err != nil {
        log.Fatal(err)
    }
    logPick.Init(logEtcdList)

    //监控etcd配置日志节点变化
    go etcd.Watch(cfg.Etcd.Key)

    select {}
}
