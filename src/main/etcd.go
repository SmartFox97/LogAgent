package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"logAgent/tailf"
	"strings"
	"time"
)

type EtcdClient struct {
	Client *clientv3.Client
	Keys []string
}

var (
	etcdCli *EtcdClient
)


func initEtcd(addr string) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdCli = &EtcdClient{
		Client: cli,
	}
	return
}

func getCollectConfig(key string) (err error) {

	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}
	for _,ipv4 := range localIPArray{
		etcdKey := fmt.Sprintf("%s%s",key,ipv4)
		etcdCli.Keys = append(etcdCli.Keys , etcdKey)
		logs.Debug("Create ETCD Key : %v" , etcdKey)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := etcdCli.Client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd failed, err:%v", err)
			continue
		}
		cancel()
		logs.Debug("resp from etcd: %v", resp)
		for _, ev := range resp.Kvs {
			if string(ev.Key) == etcdKey {
				var collects []tailf.Collect
				err = json.Unmarshal(ev.Value,&collects)
				if err != nil {
					logs.Error("unmarshal failed, err: %v", err)
					continue
				}
				logs.Debug("Get Config Success.Config：%v",collects)
				AppConfig.collects = append(AppConfig.collects , collects...)
			}
		}
	}
	//监听Key修改
	initEtcdWatcher()
	return
}

func initEtcdWatcher() {
	for _,v := range etcdCli.Keys{
		go watchKey(v)
	}
}

func watchKey(key string) (err error){
	for {
		rch := etcdCli.Client.Watch(context.Background(),key)
		var collects []tailf.Collect
		var getSuccess = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted",key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collects)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getSuccess = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getSuccess {
				logs.Debug("get config from etcd succ, %v", collects)
				err = tailf.UpdateTask(collects)
			}
		}
	}
}