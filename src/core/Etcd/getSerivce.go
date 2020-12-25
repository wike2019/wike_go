package Etcd

import (
	"encoding/json"
	"fmt"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
	"regexp"
)


func(this *EtcdCtl) LoadService (name string) ([]*ServiceInfo,error) {
	res,err:=this.Get(key_prefix,WithPrevKV())
	if err!=nil{
		return nil,err
	}
	Services:=make([]*ServiceInfo,0)

	for _,item:=range res.Kvs{

		Services=this.parseService(item.Key,item.Value,Services,name)
	}
	return Services,nil
}
func(this *EtcdCtl) parseService(key []byte,value []byte,Services []*ServiceInfo,name string) []*ServiceInfo {
	str:=fmt.Sprintf("%s([^/]+)/([^/]+)",key_prefix)
	reg:=regexp.MustCompile(str)
	if reg.Match(key){
		idandname:=reg.FindSubmatch(key)
		sname:=idandname[2]
		if string(sname)==name {
			var service ServiceInfo
			json.Unmarshal(value, &service)
			Services=append(Services,&service)
		}
	}
	return Services
}
func(this *EtcdCtl)  Seletor(Services []*ServiceInfo,selectType int,ip string)(LoadBalance.NodeBalance,error)  {
	if LB.NeedReLoad(len(Services)){
		LB.Clear()
		for _,node:=range Services{
			node.CWeight=0
			LB.AddNode(node)
		}
	}
	if selectType== LoadBalance.RoundRobinByWeight {
		return LB.RoundRobinByWeight()
	}
	if selectType==LoadBalance.RoundRobin {
		return LB.RoundRobin()
	}
	if selectType==LoadBalance.SelectByIPHash {
		return LB.SelectByIPHash(ip)
	}
	if selectType==LoadBalance.SelectByWeightRand {
		return LB.SelectByWeightRand()
	}
	if selectType==LoadBalance.SelectByRand {
		return LB.SelectByRand()
	}
	return LB.RoundRobinByWeight()
}