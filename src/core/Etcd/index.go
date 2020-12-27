package Etcd

import (
	"context"
	"github.com/wike2019/wike_go/src/core/Ioc"

	"go.etcd.io/etcd/clientv3"
)
type EtcdCtl struct {
	EtcdClient *clientv3.Client  `inject:"-"`
	ctx context.Context
}
//实现bean接口
func (this *EtcdCtl) Name() string{
	return  "EtcdCtl"
}

func init() {
	Ioc.New().Beans(New())
}
var cacheId =make(map[string]clientv3.LeaseID)
const keyPrefix ="/coreV1/services/"

func New()  *EtcdCtl{
	return &EtcdCtl{ctx:context.Background()}
}











