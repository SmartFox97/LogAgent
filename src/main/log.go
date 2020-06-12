package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func getLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "info":
		return logs.LevelInfo
	case "warn":
		return logs.LevelWarn
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug

}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = AppConfig.agentLog
	config["level"] = getLogLevel(AppConfig.logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		str := fmt.Sprintf("marshal failed, err:", err)
		err = errors.New(str)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}