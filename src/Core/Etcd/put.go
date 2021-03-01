package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)

//插入数据
func (this *EtcdCtl) Put(key string,value string,attrs ...*OperationAttr) Result.Any {
	kv:=clientv3.NewKV(this.EtcdClient)
	leaseId:=OperationAttrs(attrs).
		Find(ATTR_Lease).
		Unwrap_Or(nil)[0]
	if leaseId !=nil {
		return Result.Result(kv.Put(this.ctx,key,value,clientv3.WithLease(leaseId.(clientv3.LeaseID)))).Unwrap()[0]

	}

	time:=OperationAttrs(attrs).
		Find(ATTR_WithTime).
		Unwrap_Or([]Result.Any{0})[0].(int64)
	if time!=0{
		leaseId:= Result.Result(this.EtcdClient.Grant(this.ctx, time)).Unwrap()[0].(*clientv3.LeaseGrantResponse).ID;
		return Result.Result(kv.Put(this.ctx,key,value,clientv3.WithLease(leaseId))).Unwrap()[0]

	}
	return Result.Result(kv.Put(this.ctx,key,value)).Unwrap()[0]
}
