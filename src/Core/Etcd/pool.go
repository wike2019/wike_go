package Etcd

import (
	"github.com/wike2019/wike_go/src/Core/Bean"
	"github.com/wike2019/wike_go/src/Result"

	"sync"
)

var etcdPool *sync.Pool
//池
func init()  {
	etcdPool=&sync.Pool{
		New: func() interface{}{
			var value *EtcdCtl
			EtcdClient:= Result.Result(Bean.New().Get(value)).Unwrap()[0].(*EtcdCtl)
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

