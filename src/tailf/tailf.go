package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

//配置文件
type Collect struct {
	Topic string `json:"topic"`
	LogPath string `json:"logpath"`
}

//一条日志
type TextMsg struct {
	Topic string
	Msg string
}

type TailObj struct {
	tails *tail.Tail
	conf Collect
	status   int
	exitChan chan int
}

type TailObjMgr struct {
	tails []*TailObj
	msgChan chan *TextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func UpdateTask(configs []Collect) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, oneConf := range configs {
		var isRunning = false
		for _, obj := range tailObjMgr.tails {
			if oneConf.LogPath == obj.conf.LogPath && oneConf.Topic == obj.conf.Topic{
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}
		createTask(oneConf)
	}

	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tails {
		obj.status = StatusDelete
		for _, oneConf := range configs {
			if oneConf.LogPath == obj.conf.LogPath && oneConf.Topic == obj.conf.Topic {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}
	tailObjMgr.tails = tailObjs
	return
}

func createTask(config Collect){
	obj := &TailObj{
		conf: config,
		exitChan: make(chan int, 1),
		status: StatusNormal,
	}

	tails, err := tail.TailFile(config.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Error("collect filename[%s] failed, err:%v", config.LogPath, err)
		return
	}
	obj.tails = tails

	tailObjMgr.tails = append(tailObjMgr.tails , obj)

	go readFromTail(obj)
}

func InitTailf(config []Collect, chanSize int) (err error) {

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(config) == 0 {
		err = fmt.Errorf("invalid config for log collect, conf:%v", config)
		return
	}

	for _ , v := range config{
		createTask(v)
	}
	return
}

func readFromTail(tailObj *TailObj) {
	for true {
		select {
		case msg, ok := <- tailObj.tails.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n", tailObj.tails.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Topic: tailObj.conf.Topic,
				Msg: msg.Text,
			}
			logs.Debug("Got msg:%v",textMsg)
			tailObjMgr.msgChan <- textMsg
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.conf)
			return
		}
	}
}

func GetMessage() (msg *TextMsg) {
	msg = <- tailObjMgr.msgChan
	return
}


