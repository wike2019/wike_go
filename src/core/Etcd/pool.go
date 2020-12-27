package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/core/Ioc"

	"sync"
)

var etcdPool *sync.Pool
//池
func init()  {
	etcdPool=&sync.Pool{
		New: func() interface{}{
			var value *EtcdCtl
			EtcdClient:= Result.Result(Ioc.New().Get(value)).Unwrap().(*EtcdCtl)
			return EtcdClient
		},
	}
}
func EtcdCache() *EtcdCtl {
	return etcdPool.Get().(*EtcdCtl)
}
func ReleaseEtcdCache(cache *EtcdCtl){
	etcdPool.Put(cache)
}

