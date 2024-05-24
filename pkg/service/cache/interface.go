package cache

import (
	"time"
)

// 缓存接口
type Cache interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string, obj interface{}) bool
	Delete(k string) bool
}

// 实际缓存服务
type Service struct {
	cache Cache
}

func ServiceCache(cache Cache) *Service {
	return &Service{cache: cache}
}

// 缓存业务函数
func FindWithCallBack[T any](key string, time time.Duration, CacheService *Service, fn func() T) T {
	var data T
	ok := CacheService.cache.Get(key, &data)
	if ok {
		return data
	}
	data = fn()
	CacheService.cache.Set(key, data, time)
	return data
}
