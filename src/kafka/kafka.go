package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(ServAddr string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{ServAddr}, config)
	if err != nil {
		logs.Error("producer close, err:", err)
		return
	}
	logs.Debug("init kafka Success")
	return
}

func SendToKafka(topic ,message string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(message)
	pid, offset ,err := client.SendMessage(msg)
	if err != nil {
		logs.Debug("Send Message Failed,err:",err)
		return
	}
	logs.Debug("Send Message Succes! pid:%v offset:%v",pid,offset)
	return
}
