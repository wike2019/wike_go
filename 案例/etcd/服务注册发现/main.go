package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Etcd"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
)

func main()  {
	Ioc.New().Config(NewEtcdConfig())
	Ioc.New().ApplyAll()
	Instance:= Etcd.EtcdCache()
	defer Etcd.ReleaseEtcdCache(Instance)
	id1:= "id1"
	//注册服务
	Instance.RegService(Etcd.ServiceInfo{
		ServiceID: id1 ,
		ServiceName: "wike3",
		ServiceAddr: "192.168.3.3:8081",
		ServiceType:  "http",
		DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready},
		ServiceHost:"wike.com",

	})
	//发现服务
	m,_:=Instance.LoadService("wike3")
	fmt.Println(Instance.Seletor(m,LoadBalance.RoundRobinByWeight,"192.168.127.1"))
	//反注册
	Instance.UnregService(id1)
}
