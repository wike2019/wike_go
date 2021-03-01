package LoadBalance

import (
	"fmt"
	"sort"
)

//平滑加权轮询
func(this *LoadBalance) RoundRobinByWeight() (NodeBalance,error) {
	ok:=this.allDown()
	if ok {
		return nil,fmt.Errorf("节点都不可用")
	}
	for _,s:=range this.nodes{
		s.SetCWeight(s.GetWeight()+s.GetCWeight())
	}
	sort.Sort(this.nodes)
	max:=this.nodes[0] //返回最大 作为命中服务
	SumWeight:=0
	for _,server:=range this.nodes{
		SumWeight=SumWeight+server.GetWeight()  //计算加权总和
	}
	max.SetCWeight(max.GetCWeight()-SumWeight)

	if max.GetStatus()!=Ready {

		return  this.RoundRobinByWeight()
	}
	return max,nil
}