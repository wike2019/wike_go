package cache

import (
	"github.com/wike2019/wike_go/src/core/ioc"
	"github.com/wike2019/wike_go/src/core/redis"
	"github.com/wike2019/wike_go/src/result"
	"sync"
	"time"
)

var NewsCachePool *sync.Pool

func init()  {
	NewsCachePool=&sync.Pool{
		New: func() interface{} {
			RedisStringOperation:= result.Result(ioc.New().ExprData["RedisStringOperation"]).Unwrap()
			return redis.NewSimpleCache(RedisStringOperation.(*redis.RedisStringOperation),time.Second*15,redis.Serilizer_JSON)
		},
	}
}
func RedisCache() *redis.SimpleCache {
	return NewsCachePool.Get().(*redis.SimpleCache)
}
func ReleaseRedisCache(cache *redis.SimpleCache){
	NewsCachePool.Put(cache)
}

