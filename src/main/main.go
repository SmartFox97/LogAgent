package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logAgent/kafka"
	"logAgent/tailf"
)

func main() {
	fileName := "./conf/config.cfg"
	confType := "ini"

	//读取配置文件
	err := loadConfig(confType, fileName)
	if err != nil {
		fmt.Println("Load Config Failed, err:", err)
		panic("load conf failed")
		return
	}

	err = initLogger()
	if err != nil {
		fmt.Printf("load logger failed, err:%v\n", err)
		panic("load logger failed")
		return
	}
	logs.Debug("load conf succ, config:%v", AppConfig)

	//判断是否需要加载etcd
	if AppConfig.etcdenable {
		err = initEtcd(AppConfig.etcdAddr)
		if err != nil {
			logs.Error("Init ETCD Failed , err: %v",err)
			return
		}
		err = getCollectConfig(AppConfig.etcdKey)
		if err != nil {
			logs.Error("Get CollectConfig Failed err: %v" ,err)
			return
		}
	}

	err = tailf.InitTailf(AppConfig.collects,AppConfig.chanSize)
	if err != nil {
		logs.Error("Init Tail Failed, err : %v",err)
		return
	}

	err = kafka.InitKafka(AppConfig.kafkaSer)
	if err != nil {
		logs.Error("init Kafka failed, err:%v", err)
		return
	}
	logs.Debug("System Initialize Success!")

	//加载完成,启动监控逻辑
	err = serverRun()
	if err != nil {
		logs.Error("Server Run Failed, err : %v",err)
		return
	}

	logs.Info("Program Exited")

}
