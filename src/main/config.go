package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"logAgent/tailf"
)

var (
	AppConfig *Config
)

type Config struct {
	logLevel string
	agentLog string
	chanSize int
	kafkaSer string

	collects []tailf.Collect
}

func loadCollectConfig(conf config.Configer) (err error) {
	var collects tailf.Collect
	collects.LogPath = conf.String("collect::LogPath")
	if len(collects.LogPath) == 0 {
		err = errors.New("invalid collect::log_path")
		return
	}
	collects.Topic = conf.String("collect::Topic")
	if len(collects.Topic) == 0 {
		err = errors.New("invalid collect::topic")
	}
	AppConfig.collects = append(AppConfig.collects,collects)
	return
}

func loadConfig(confType string,filename string)(err error){
	//读取配置文件
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		return
	}
	AppConfig = &Config{}

	AppConfig.logLevel = conf.String("logAgent::LogLevel")
	if len(AppConfig.logLevel) == 0 {
		AppConfig.logLevel = "debug"
	}

	AppConfig.agentLog = conf.String("logAgent::AgentLog")
	if len(AppConfig.agentLog) == 0 {
		AppConfig.agentLog = "./logs"
	}

	AppConfig.chanSize,err = conf.Int("logAgent::ChanSize")
	if err != nil {
		AppConfig.chanSize = 100
	}

	kafkaServer := conf.String("kafka::Server")
	if len(kafkaServer) == 0 {
		err = fmt.Errorf("invalid kafka Server addr")
	}
	kafkaPort,err:= conf.Int("kafka::Port")
	if err != nil {
		return
	}
	AppConfig.kafkaSer = fmt.Sprintf("%s:%d",kafkaServer,kafkaPort)

	err = loadCollectConfig(conf)
	if err != nil {
		fmt.Printf("load collect conf failed, err:%v\n", err)
		return
	}
	return
}