package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/config"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	fileName := "./conf/config.cfg"
	confType := "ini"
	conf, err := config.NewConfig(confType, fileName)
	if err != nil {
		return
	}
	EtcdSerAddr := conf.String("etcd::SerAddr")
	key := "/LogAgent/config/192.168.100.131"
	GetConfig(EtcdSerAddr,key)
}

func GetConfig(addr string,key string)  {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
