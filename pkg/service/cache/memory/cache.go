package memory

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

// 内存缓存
type Cache struct {
	MemoryCache *cache.Cache
}

func NewCache() *Cache {
	cache := cache.New(5*time.Minute, 10*time.Minute)
	return &Cache{MemoryCache: cache}
}

// 接口方法
func (this *Cache) Get(k string, obj interface{}) bool {
	data, found := this.MemoryCache.Get(k)
	if found {
		err := json.Unmarshal(data.([]byte), obj)
		if err != nil {
			return false
		}
		return found
	}
	return false
}
func (this *Cache) Set(k string, x interface{}, d time.Duration) {
	b, _ := json.Marshal(x)
	this.MemoryCache.Set(k, b, d)
}

func (this *Cache) Delete(k string) bool {
	this.MemoryCache.Delete(k)
	return true
}
