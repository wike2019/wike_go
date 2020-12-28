package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
)

type Target struct {
	TargetId   string
	*LoadBalance.DefaultBalance
}
func  (this *Target) GetNode() interface{}  {
	return this
}


var LB= LoadBalance.New()



func main()  {
	Services:=make([]LoadBalance.NodeBalance,0)
	Services=append(Services,
					&Target{TargetId:"id1" ,DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:"id1"}},
					&Target{TargetId:"id2" ,DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:"id2"}},
					&Target{TargetId:"id3" ,DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:"id3"}},
	)
	LB.AddNode(Services...)

	selectType:=LoadBalance.RoundRobinByWeight
	if selectType== LoadBalance.RoundRobinByWeight {
		 fmt.Println(LB.RoundRobinByWeight())
	}
	if selectType==LoadBalance.RoundRobin {
		fmt.Println(LB.RoundRobin())
	}
	if selectType==LoadBalance.SelectByIPHash {
		fmt.Println( LB.SelectByIPHash("192.168.3.2"))
	}
	if selectType==LoadBalance.SelectByWeightRand {
		fmt.Println( LB.SelectByWeightRand())
	}
	if selectType==LoadBalance.SelectByRand {
		fmt.Println( LB.SelectByRand())
	}
	if LB.NeedReLoad(Services){
		fmt.Println("没有修改过")
	}
	Services=append(Services,&Target{TargetId:"id4" ,DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:"id4"}})
	if LB.NeedReLoad(Services){
		fmt.Println("修改过")
	}

	node,_:=LB.RoundRobinByWeight()
	fmt.Println(node.GetNode()) //节点信息
}
