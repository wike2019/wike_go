package LoadBalance

type NodeBalance interface {
	GetNode() interface{}
	GetCWeight() int
	GetWeight() int
	GetStatus() string
	SetCWeight(val int)
	SetWeight(Weight int)
	GetId() string
}

