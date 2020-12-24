package etcd

import (
	"encoding/json"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
)

type ServiceInfo struct {
	ServiceID string
	ServiceName string
	ServiceAddr string
	ServiceType string
	ServiceHost string
	ServiceWight int
	Status string
	CWeight int
}

func  (this *ServiceInfo) GetCWeight() int {
	return this.CWeight
}
func  (this *ServiceInfo) SetCWeight(val int)  {
	this.CWeight=val
}
func  (this *ServiceInfo) GetStatus() string {
	return  this.Status
}

func  (this *ServiceInfo) GetNode() interface{}  {
	return this
}
func  (this *ServiceInfo) GetWeight() int {
	return this.ServiceWight
}

var LB= LoadBalance.New()
//注册服务
func(this *EtcdCtl) RegService(info ServiceInfo){
	LB.AddNode(&info)
	key:=key_prefix+info.ServiceID+"/"+info.ServiceName
	reid:=this.Lease(30)
	b, _ := json.Marshal(info)
	this.Put(key,string(b),WithLease(reid))
	cachId[info.ServiceID]=reid
	this.KeepAlive(reid)
}

//反注册服务
func(this *EtcdCtl) UnregService(id string) error  {
	this.DelLease(cachId[id])
	_,err:=this.Del(key_prefix+id,WithPrevKV())
	return err
}
