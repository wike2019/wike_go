package Etcd

import (
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/core/Ioc"

	"sync"
)

var etcdPool *sync.Pool

func init()  {
	etcdPool=&sync.Pool{
		New: func() interface{}{
			EtcdClient:= Result.Result(Ioc.New().ExprData["EtcdCtl"]).Unwrap().(*EtcdCtl)
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

