package LoadBalance



type NodeBalanceSlice  []NodeBalance
func (p NodeBalanceSlice) Len() int           { return len(p) }
func (p NodeBalanceSlice) Less(i, j int) bool { return p[i].GetCWeight() > p[j].GetCWeight() } //从大到小排序
func (p NodeBalanceSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }


type LoadBalance struct { //负载均衡类
	nodes NodeBalanceSlice
	curIndex int  //只想当前访问的index

}


func(this *LoadBalance) AddNode(nodeSlice ... NodeBalance)  {
	for _,node:=range nodeSlice{
		this.nodes=append(this.nodes,node)
	}
}
