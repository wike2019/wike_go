package LoadBalance

import (
	"sort"
	"strings"
)

func New() *LoadBalance {
	return &LoadBalance{nodes:make(NodeBalanceSlice,0)}
}

func(this *LoadBalance) NeedReLoad(service NodeBalanceSlice) bool { //判断是否要重新计算权重
	if  this.nodes.Len() ==len(service){
		item:=service
		sort.SliceStable(item, func(i, j int) bool {
			return  strings.Compare(item[i].GetId(),item[j].GetId())>0
		})
		sort.SliceStable(this.nodes, func(i, j int) bool {
			return  strings.Compare(item[i].GetId(),item[j].GetId())>0
		})
		for i:=0;i<this.nodes.Len() ; i++ {
			if strings.Compare(item[i].GetId(),this.nodes[i].GetId())!=0{
				return  true
			}
		}
		return  false
	}
	return true
}
func(this *LoadBalance) Clear() { //随机算法
	this.nodes=make(NodeBalanceSlice,0)
}
type DefaultBalance struct {
	CWeight int
	Status string
	Weight int
	Id string
}


func  (this *DefaultBalance) GetCWeight() int {
	return this.CWeight
}
func  (this *DefaultBalance) GetId() string {
	return this.Id
}
func  (this *DefaultBalance) GetWeight() int {
	return this.Weight
}
func  (this *DefaultBalance) SetCWeight(val int)  {
	this.CWeight=val
}
func  (this *DefaultBalance) GetStatus() string {
	return  this.Status
}
func  (this *DefaultBalance) SetWeight(Weight int )  {
	  this.Weight=Weight
}
