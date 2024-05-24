package bloomfilter

import (
	"github.com/hugh2632/bloomfilter"
)

// 布隆过滤器
type Bloom struct {
	Bloom bloomfilter.IFilter
}

// 默认布隆过滤器
func NewBloom() *Bloom {
	return &Bloom{Bloom: bloomfilter.NewMemoryFilter(make([]byte, 10240), bloomfilter.DefaultHash...)}
}

func (this *Bloom) Exists(key string) bool {
	return this.Bloom.Exists([]byte(key))
}
func (this *Bloom) Clear() error {
	return this.Bloom.Clear()
}
func (this *Bloom) Add(key string) {
	this.Bloom.Push([]byte(key))
}
