package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

// redis 缓存
type Cache struct {
	redis *redis.Client
}

const (
	NoExpiration      time.Duration = 0
	DefaultExpiration time.Duration = 5 * time.Minute
)

func NewCache(redisCache *redis.Client) *Cache {
	return &Cache{
		redis: redisCache,
	}
}

// 接口方法
func (this *Cache) Set(k string, x interface{}, d time.Duration) {
	b, _ := json.Marshal(x)
	this.redis.Set(
		context.Background(),
		k,
		b,
		d,
	)
}
func (this *Cache) Get(k string, obj interface{}) bool {
	b, err := this.redis.Get(context.Background(), k).Result()
	if err == nil {
		err = json.Unmarshal([]byte(b), obj)
		if err == nil {
			return true
		}
	}
	return false
}

func (this *Cache) Delete(k string) bool {
	err := this.redis.Del(context.Background(), k).Err()
	if err == nil {
		return true
	}
	return false
}
