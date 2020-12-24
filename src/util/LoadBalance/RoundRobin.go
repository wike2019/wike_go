package LoadBalance

import "fmt"

func(this *LoadBalance) RoundRobin() (NodeBalance,error)  {
	if this.allDown() {
		return nil,fmt.Errorf("节点都不可用")
	}
	node:=this.nodes[this.curIndex]

	this.curIndex=(this.curIndex+1) % len(this.nodes)
	if node.GetStatus()!=Ready {
		return  this.RoundRobin()
	}
	return  node,nil
}