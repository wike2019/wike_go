package main

import (
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_app/controller"
	redisInit "github.com/wike2019/wike_app/redis"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/cache"
	"github.com/wike2019/wike_go/v2/func/cache/memory"
	"github.com/wike2019/wike_go/v2/func/cache/redis"
	"github.com/wike2019/wike_go/v2/fxTags"
)

func main() {
	//一个最简单的例子
	g := core.God()
	g.GlobalUse(core.CustomRecover(g))
	g.Provide(fxTags.CreateInterFace(memory.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_memory"), nil))
	g.Provide(fxTags.CreateInterFace(redis.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_redis"), nil))
	g.Provide(config.Config, redisInit.InitRedis).
		Provide(fxTags.Create(cache.ServiceCache, "", fxTags.ParamList(fxTags.CreateTag("cache_redis")))). //选择redis作为缓存服务的存储

		Mount(controller.NewRouter3,
			fxTags.ParamList(
				fxTags.CreateTag("cache_memory"), fxTags.CreateTag("cache_redis"))).Run()
}
