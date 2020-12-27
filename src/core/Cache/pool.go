package Cache

import (
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Redis"
	"sync"
	"time"
)

var NewsCachePool *sync.Pool
//缓冲池
func init()  {
	NewsCachePool=&sync.Pool{
		New: func() interface{} {
			var value *Redis.RedisStringOperation
			RedisStringOperation:= Result.Result(Ioc.New().Get(value)).Unwrap()
			//缓存5分钟
			return Redis.NewSimpleCache(RedisStringOperation.(*Redis.RedisStringOperation),time.Second*300, Redis.Serilizer_JSON)
		},
	}
}
func RedisCache() *Redis.SimpleCache {
	return NewsCachePool.Get().(*Redis.SimpleCache)
}
func ReleaseRedisCache(cache *Redis.SimpleCache){
	NewsCachePool.Put(cache)
}

