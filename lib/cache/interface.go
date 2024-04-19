package cache

import (
	"github.com/wike2019/wike-go/lib/bloom"
	"time"
)

// 缓存接口
type Cache interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string, obj interface{}) bool
}

// 实际缓存服务
type Service struct {
	cache Cache
	bloom *bloom.Bloom
}

func ServiceCache(cache Cache, bloom *bloom.Bloom) *Service {
	return &Service{cache: cache, bloom: bloom}
}

// 缓存业务函数
func FindWithCallBack[T any](key string, time time.Duration, CacheService *Service, fn func() T) T {
	var data T
	isExists := CacheService.bloom.Bloom.Exists([]byte(key))
	if isExists {
		ok := CacheService.cache.Get(key, &data)
		if ok {
			return data
		}
		data = fn()
		CacheService.cache.Set(key, data, time)
		return data
	}
	CacheService.bloom.Bloom.Push([]byte(key))
	data = fn()
	CacheService.cache.Set(key, data, time)
	return data
}
