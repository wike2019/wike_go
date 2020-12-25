package Etcd

import (
	"go.etcd.io/etcd/clientv3"
)

func (this *EtcdCtl) Get(key string,attrs ...*OperationAttr) (*clientv3.GetResponse,error) {
	kv:=clientv3.NewKV(this.EtcdClient)
	prev:=OperationAttrs(attrs).
		Find(ATTR_WithPrevKV).
		Unwrap_Or(nil)
	if prev !=nil {
		return  kv.Get(this.ctx,key,clientv3.WithPrefix())
	}

	return  kv.Get(this.ctx,key)
}