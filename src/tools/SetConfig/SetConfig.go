package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"go.etcd.io/etcd/clientv3"
	"logAgent/tailf"
	"time"
)


func SetLogConfToEtcd(etcdAddr string,etcdKey string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}
	fmt.Println("connect succ")
	defer cli.Close()
	var logConfArr []tailf.Collect
	logConfArr = append(
		logConfArr,
			tailf.Collect{
				LogPath: "D:\\LogAgent.log",
				Topic:   "MyTest",
			},
		)
	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}



	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, etcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	cli.Delete(ctx, etcdKey)
	cancel()
	return

}

func main() {
	fileName := "./conf/config.cfg"
	confType := "ini"
	conf, err := config.NewConfig(confType, fileName)
	if err != nil {
		return
	}
	EtcdSerAddr := conf.String("etcd::SerAddr")
	EtcdKey := conf.String("etcd::ETCDKey")
	ipv4 := "192.168.100.131"
	etcdKey := fmt.Sprintf("%s%s",EtcdKey,ipv4)
	SetLogConfToEtcd(EtcdSerAddr,etcdKey)
}
