package Etcd

import (
	"context"
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)


//续租
func (this * EtcdCtl)  KeepAlive(LeaseID clientv3.LeaseID) {
	keepaliveRes:=Result.Result(this.EtcdClient.KeepAlive(this.ctx,LeaseID)).Unwrap()[0].(<-chan *clientv3.LeaseKeepAliveResponse )
	go lisKeepAlive(keepaliveRes)

}

func (this * EtcdCtl)  KeepAliveWithContext(ctx context.Context,LeaseID clientv3.LeaseID) {
	keepaliveRes:=Result.Result(this.EtcdClient.KeepAlive(ctx,LeaseID)).Unwrap()[0].(<-chan *clientv3.LeaseKeepAliveResponse )
	go lisKeepAlive(keepaliveRes)

}
func lisKeepAlive(keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse)  {
	for{
		select {
		case ret:=<-keepaliveRes:
			if ret!=nil{
			}
		}
	}
}
