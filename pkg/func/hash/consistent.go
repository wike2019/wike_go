package hash

import "github.com/stathat/consistent"

type Hash struct {
	*consistent.Consistent
}

func New() *Hash {
	return &Hash{
		Consistent: consistent.New(),
	}
}
func (this *Hash) Get(key string) (string, error) {
	return this.Consistent.Get(key)
}
func (this *Hash) Add(node string) {
	this.Consistent.Add(node)
}
