package etcd

import "go.etcd.io/etcd/clientv3"

func (this *EtcdCtl) Del(key string,attrs ...*OperationAttr) (*clientv3.DeleteResponse,error) {
	kv:=clientv3.NewKV(this.EtcdClient)
	prev:=OperationAttrs(attrs).
		Find(ATTR_WithPrevKV).
		Unwrap_Or(nil)
	if prev !=nil {
		return  kv.Delete(this.ctx,key,clientv3.WithPrefix())
	}

	return  kv.Delete(this.ctx,key)
}

func (this *EtcdCtl) DelLease(key clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	lease := clientv3.NewLease(this.EtcdClient)
	return  lease.Revoke(this.ctx,key)
}
