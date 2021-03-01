package Etcd

import (
	"context"
	"github.com/wike2019/wike_go/src/Core/Bean"
	"go.etcd.io/etcd/clientv3"
)

var cacheId =make(map[string]clientv3.LeaseID)
const keyPrefix ="/coreV1/services/"


type EtcdCtl struct {
	EtcdClient *clientv3.Client  `inject:"-"`
	ctx context.Context
}
//实现bean接口
func (this *EtcdCtl) Name() string{
	return  "EtcdCtl"
}
func init() {
	Bean.New().Beans(New())
}

func New()  *EtcdCtl{
	return &EtcdCtl{ctx:context.Background()}
}











