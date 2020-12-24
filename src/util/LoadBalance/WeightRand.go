package LoadBalance

import (
	"fmt"
	"math/rand"
	"time"
)

func(this *LoadBalance) SelectByWeightRand() (NodeBalance,error) { //加权随机算法(改良算法)
	if this.allDown() {
		return nil,fmt.Errorf("节点都不可用")
	}
	rand.Seed(time.Now().UnixNano())
	sumList:=make([]int,len(this.nodes))
	sum:=0
	for	 i,_:= range  this.nodes{
		temp:=this.nodes[i].(NodeBalance).GetWeight()
		sum+=temp
		sumList[i]=sum
	}
	rad:=rand.Intn(sum) //[)
	for index,value:=range sumList{
		if rad<value{
			if this.nodes[index].GetStatus()!=Ready{
				return  this.SelectByWeightRand()
			}
			return  this.nodes[index],nil
		}
	}
	return nil,fmt.Errorf("数据异常")
}