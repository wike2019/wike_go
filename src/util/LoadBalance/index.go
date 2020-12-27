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
type DefaultBalance struct {
	CWeight int
	Status string
	Weight int
}


func  (this *DefaultBalance) GetCWeight() int {
	return this.CWeight
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
