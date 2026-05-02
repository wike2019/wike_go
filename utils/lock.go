package utils

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLock struct {
	client   *redis.Client
	mutex    *redislock.Client
	LockData *redislock.Lock
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}
func (RedisLock) String() string {
	return "redis"
}

func (r *RedisLock) Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error) {
	if r.mutex == nil {
		r.mutex = redislock.New(r.client)
	}
	lock, err := r.mutex.Obtain(context.TODO(), key, time.Duration(ttl)*time.Second, options)
	r.LockData = lock
	// 设置一个定时任务来自动续期锁
	go func() {
		for {
			time.Sleep(5 * time.Second)                                              // 每 5 秒续期一次
			err := r.LockData.Refresh(context.Background(), 10*time.Second, options) // 延长锁的有效期
			if err != nil {
				return
			}
		}
	}()
	return r.LockData, err
}
