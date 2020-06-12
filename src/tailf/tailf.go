package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"time"
)

//配置文件
type Collect struct {
	Topic string
	LogPath string
}

//一条日志
type TextMsg struct {
	Topic string
	Msg string
}

type TailObj struct {
	tails *tail.Tail
	conf Collect
}

type TailObjMgr struct {
	tails []*TailObj
	msgChan chan *TextMsg
}

var (
	tailObjMgr *TailObjMgr
)

func createTask(config Collect){
	obj := &TailObj{
		conf: config,
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
		msg, ok := <- tailObj.tails.Lines
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
	}
}

func GetMessage() (msg *TextMsg) {
	msg = <- tailObjMgr.msgChan
	return
}


