package main

import (
	"github.com/astaxie/beego/logs"
	"logAgent/kafka"
	"logAgent/tailf"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetMessage()
		err = kafka.SendToKafka(msg.Topic,msg.Msg)
		if err != nil {
			logs.Error("send to kafka failed, err:%v ,Will try After 10s", err)
			time.Sleep(10*time.Second)
			continue
		}
	}
	return
}