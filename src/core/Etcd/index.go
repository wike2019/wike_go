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

func init() {
	Ioc.New().Beans(New())
}
var cachId=make(map[string]clientv3.LeaseID)
var key_prefix ="/coreV1/services/"

func New()  *EtcdCtl{
	return &EtcdCtl{ctx:context.Background()}
}
func (this *EtcdCtl) Name() string{
	return  "EtcdCtl"
}










