package main

import (
	"fmt"

)

type EtcdCall struct {
	Etcdctl *Etcd.EtcdCtl   `inject:"-"`
}

func init() {
	Ioc.New().Beans(new(EtcdCall))
}

func(this *EtcdCall) Name() string  {
	return "EtcdCall"
}
func(this *EtcdCall) Get()   {
	fmt.Println(this.Etcdctl.Get("key"))
}
func(this *EtcdCall) GetWithPrevKV()   {
	fmt.Println(this.Etcdctl.Get("key",Etcd.WithPrevKV()))
}
func(this *EtcdCall) SetWithTimeOut()   {
	fmt.Println(this.Etcdctl.Put("key","age",Etcd.WithTime(10)))
}
func(this *EtcdCall) Set()   {
	fmt.Println(this.Etcdctl.Put("key","age"))
}
func(this *EtcdCall) SetLeaseId()   {
	fmt.Println(this.Etcdctl.Put("key","age",Etcd.WithLease(this.Etcdctl.Lease(10))))
}
func(this *EtcdCall) Del()   {
	fmt.Println(this.Etcdctl.Del("key"))
}
func(this *EtcdCall) DelWithPrevKV()   {
	fmt.Println(this.Etcdctl.Del("key",Etcd.WithPrevKV()))
}
func(this *EtcdCall) DelLease()   {
	id:=this.Etcdctl.Lease(10)
	fmt.Println(this.Etcdctl.DelLease(id))
}
func(this *EtcdCall) KeepAlive()   {
	id:=this.Etcdctl.Lease(10)
	fmt.Println(this.Etcdctl.KeepAlive(id))
}