package Etcd

import (
	"fmt"
	"github.com/wike2019/wike_go/src/result"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func (this *EtcdCtl) Put(key string,value string,attrs ...*OperationAttr) interface{} {
	kv:=clientv3.NewKV(this.EtcdClient)
	leaseId:=OperationAttrs(attrs).
		Find(ATTR_Lease).
		Unwrap_Or(nil)
	fmt.Println(leaseId)
	if leaseId !=nil {
		fmt.Println(3333333333333)
		return Result.Result(kv.Put(this.ctx,key,value,clientv3.WithLease(leaseId.(clientv3.LeaseID)))).Unwrap()

	}
	time:=OperationAttrs(attrs).
		Find(ATTR_WithTime).
		Unwrap_Or(time.Second*0).(int64)
	if time!=0{
		leaseId:= Result.Result(this.EtcdClient.Grant(this.ctx, time)).Unwrap().(*clientv3.LeaseGrantResponse).ID;
		return Result.Result(kv.Put(this.ctx,key,value,clientv3.WithLease(leaseId))).Unwrap()

	}
	return Result.Result(kv.Put(this.ctx,key,value)).Unwrap()
}
