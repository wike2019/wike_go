package LoadBalance

import (
	"fmt"
	"math/rand"
	"time"
)

func(this *LoadBalance) SelectByRand() (NodeBalance,error) { //随机算法
	if this.allDown() {
		return nil,fmt.Errorf("节点都不可用")
	}

	rand.Seed(time.Now().UnixNano())
	index:=rand.Intn(len(this.nodes))
	if this.nodes[index].GetStatus()!=Ready {
		return  this.SelectByRand()
	}
	return  this.nodes[index],nil
}