package etcd

import (
	"github.com/wike2019/wike_go/src/result"
	"go.etcd.io/etcd/clientv3"
)

func(this *EtcdCtl) Lease(time int64)clientv3.LeaseID {
	return result.Result(this.EtcdClient.Grant(this.ctx, time)).Unwrap().(*clientv3.LeaseGrantResponse).ID
}