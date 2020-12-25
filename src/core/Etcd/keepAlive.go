package Etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func (this * EtcdCtl)  KeepAlive(LeaseID clientv3.LeaseID) error{
	keepaliveRes,err:=this.EtcdClient.KeepAlive(this.ctx,LeaseID)
	if err!=nil{
		return err
	}
	go lisKeepAlive(keepaliveRes)
	return nil
}
func (this * EtcdCtl)  KeepAliveWithContext(ctx context.Context,LeaseID clientv3.LeaseID) error{
	keepaliveRes,err:=this.EtcdClient.KeepAlive(ctx,LeaseID)
	if err!=nil{
		return err
	}
	go lisKeepAlive(keepaliveRes)
	return nil
}
func lisKeepAlive(keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse)  {
	for{
		select {
		case ret:=<-keepaliveRes:
			if ret!=nil{
				fmt.Println("续租成功",time.Now())
			}
		}
	}
}
