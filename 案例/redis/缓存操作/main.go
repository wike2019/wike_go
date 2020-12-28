package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Cache"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Redis"
	"time"
)

func main()  {
	Ioc.New().Config(NewRedisConfig())
	Ioc.New().ApplyAll()
	id:="1"
	newsCache:= Cache.RedisCache()
	defer Cache.ReleaseRedisCache(newsCache)

	newsCache.DBGetter= NewUserGetter(id) //一旦缓存没有，则需要从哪去取
	newsCache.SetCahPolicy(Redis.NewCrossPolicy("tested\\d+",time.Second*60))
	//NewCrossPolicy 第一个参数是key的正则校验 ,第二个参数是 如果找不到值的缓存时间可以比正常找到的设置短防止缓存穿透
	i:=1;
	for {
		if i>5{
			break
		}
		str:=new(string)
		// 3、取缓存输出（一旦没有，上面的DBGetter会被调用)
		newsCache.GetCacheForObject("tested"+id,str)
		fmt.Println(*str)
		i++
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("key不合法")
			}
		}()
		id="woiksad"
		str:=new(string)
		newsCache.GetCacheForObject("tested"+id,str) //此时会panic 因为不符合"tested\\d+"正则
		fmt.Println(*str)
	}()


	for {
		time.Sleep(1*time.Second)
	}
}
