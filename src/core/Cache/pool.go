package Cache

import (
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Redis"
	"sync"
	"time"
)

var NewsCachePool *sync.Pool

func init()  {
	NewsCachePool=&sync.Pool{
		New: func() interface{} {
			RedisStringOperation:= Result.Result(Ioc.New().ExprData["RedisStringOperation"]).Unwrap()
			return Redis.NewSimpleCache(RedisStringOperation.(*Redis.RedisStringOperation),time.Second*15, Redis.Serilizer_JSON)
		},
	}
}
func RedisCache() *Redis.SimpleCache {
	return NewsCachePool.Get().(*Redis.SimpleCache)
}
func ReleaseRedisCache(cache *Redis.SimpleCache){
	NewsCachePool.Put(cache)
}

