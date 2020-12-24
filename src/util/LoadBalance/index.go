package LoadBalance

func New() *LoadBalance {
	return &LoadBalance{nodes:make(NodeBalanceSlice,0)}
}

func(this *LoadBalance) NeedReLoad(len int)bool { //随机算法
	return  this.nodes.Len() !=len
}
func(this *LoadBalance) Clear() { //随机算法
	this.nodes=make(NodeBalanceSlice,0)
}