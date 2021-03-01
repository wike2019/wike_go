package Etcd

import (
	"encoding/json"
	"fmt"
	"github.com/wike2019/wike_go/src/Util/LoadBalance"
	"go.etcd.io/etcd/clientv3"
	"regexp"
)

//取得数据
func(this *EtcdCtl) LoadService (name string) []LoadBalance.NodeBalance {
	res:=this.Get(keyPrefix,WithPrevKV()).Unwrap()[0].(*clientv3.GetResponse)
	Services:=make([]LoadBalance.NodeBalance,0)
	for _,item:=range res.Kvs{
		Services=this.parseService(item.Key,item.Value,Services,name)
	}
	return Services
}
//解析服务
func(this *EtcdCtl) parseService(key []byte,value []byte,Services []LoadBalance.NodeBalance,name string) []LoadBalance.NodeBalance {
	str:=fmt.Sprintf("%s([^/]+)/([^/]+)", keyPrefix)
	reg:=regexp.MustCompile(str)
	if reg.Match(key){
		id:=reg.FindSubmatch(key)
		sname:=id[2]
		if string(sname)==name {
			var service ServiceInfo
			json.Unmarshal(value, &service)
			Services=append(Services, &service)
		}
	}
	return Services
}

//选择器
func(this *EtcdCtl)  Seletor(Services LoadBalance.NodeBalanceSlice,selectType int,ip string)(LoadBalance.NodeBalance,error)  {

	if LB.NeedReLoad(Services){
		LB.Clear()
		for _,node:=range Services{
			node.SetCWeight(0)
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