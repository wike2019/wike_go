package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)


//取值
func (this *EtcdCtl) Get(key string,attrs ...*OperationAttr) *Result.ErrorResult {
	kv:=clientv3.NewKV(this.EtcdClient)
	prev:=OperationAttrs(attrs).
		Find(ATTR_WithPrevKV).
		Unwrap_Or(nil)[0]
	if prev !=nil {
		return  Result.Result(kv.Get(this.ctx,key,clientv3.WithPrefix()))
	}

	return  Result.Result(kv.Get(this.ctx,key))
}