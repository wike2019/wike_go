package Etcd

import (
	"encoding/json"
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/Util/LoadBalance"
)

type ServiceInfo struct {
	ServiceID   string
	ServiceName string
	ServiceAddr string
	ServiceType string
	ServiceHost string
    *LoadBalance.DefaultBalance
}




func  (this *ServiceInfo) GetNode() interface{}  {
	return this
}


var LB= LoadBalance.New()
//注册服务
func(this *EtcdCtl) RegService(info ServiceInfo){
	LB.AddNode(&info)
	key:= keyPrefix +info.ServiceID+"/"+info.ServiceName
	reid:=this.Lease(30)
	b, _ := json.Marshal(info)
	this.Put(key,string(b),WithLease(reid))
	cacheId[info.ServiceID]=reid
	this.KeepAlive(reid)
}

//反注册服务
func(this *EtcdCtl) UnregService(id string)   {
	this.DelLease(cacheId[id])
	Result.Result(this.Del(keyPrefix+id,WithPrevKV())).Unwrap()
}
