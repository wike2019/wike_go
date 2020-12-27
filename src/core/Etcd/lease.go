package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)
//创建租约
func(this *EtcdCtl) Lease(time int64)clientv3.LeaseID {
	return Result.Result(this.EtcdClient.Grant(this.ctx, time)).Unwrap().(*clientv3.LeaseGrantResponse).ID
}