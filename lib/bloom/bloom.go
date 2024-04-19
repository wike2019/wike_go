package bloom

import (
	"github.com/hugh2632/bloomfilter"
)

// 布隆过滤器
type Bloom struct {
	Bloom bloomfilter.IFilter
}

// 默认布隆过滤器
func DefaultBloom() *Bloom {
	return &Bloom{Bloom: bloomfilter.NewMemoryFilter(make([]byte, 10240), bloomfilter.DefaultHash...)}
}

var Clear chan struct{}

func init() {
	Clear = make(chan struct{})
}
