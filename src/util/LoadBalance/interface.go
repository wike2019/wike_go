package LoadBalance

type NodeBalance interface {
	GetNode() interface{}
	GetWeight() int
	GetCWeight() int
	SetCWeight(int)
	GetStatus() string
}