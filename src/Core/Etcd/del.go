package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)
//删除操作
func (this *EtcdCtl) Del(key string,attrs ...*OperationAttr) *Result.ErrorResult {
	kv:=clientv3.NewKV(this.EtcdClient)
	prev:=OperationAttrs(attrs).
		Find(ATTR_WithPrevKV).
		Unwrap_Or(nil)[0]
	if prev !=nil {
		return  Result.Result(kv.Delete(this.ctx,key,clientv3.WithPrefix()))
	}

	return Result.Result(kv.Delete(this.ctx,key))
}
//删除租约
func (this *EtcdCtl) DelLease(key clientv3.LeaseID) *Result.ErrorResult {
	lease := clientv3.NewLease(this.EtcdClient)
	return Result.Result(lease.Revoke(this.ctx,key))
}
