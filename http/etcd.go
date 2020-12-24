package http

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)



type EtcdConfig struct {
}
func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

func(this *EtcdConfig) Etcd() *clientv3.Client{
	config:=clientv3.Config{
		Endpoints:[]string{"192.168.3.2:32379","192.168.3.2:32380","192.168.3.2:32381"},
		DialTimeout:10*time.Second,
	}
	client,_:=clientv3.New(config)

	return client
}