package main

import (
	"fmt"

	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_app/controller"
	"github.com/wike2019/wike_app/model"
	"github.com/wike2019/wike_app/provide"
	redisInit "github.com/wike2019/wike_app/redis"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/cache"
	"github.com/wike2019/wike_go/v2/func/cache/memory"
	"github.com/wike2019/wike_go/v2/func/cache/redis"
	"github.com/wike2019/wike_go/v2/fxTags"
)

func main() {
	CatSupply := &model.CatSupply{Name: "先初始化完成再注入"}

	//一个最简单的例子
	g := core.God()
	g.Config(&model.NoParam{}).Provide(config.Config, redisInit.InitRedis).

		//上面代码还是报错，原因还是在同时有2个Animal对象依赖注入框架不知道是哪个所以需要打tag
		//.Provide(ProvideExampleSecond)
		//.Provide(ProvideExampleFirst)
		Provide(fxTags.Create(provide.ProvideExampleSecond, fxTags.CreateTag("dogSecond"), nil)).
		Provide(fxTags.Create(provide.ProvideExampleFirst, fxTags.CreateTag("dogFirst"), nil)).
		MountWithEmpty(controller.NewRouter).
		Supply(CatSupply).
		Invokes(func(r *model.CatSupply) {
			fmt.Println(r.Name, "这个是立即执行的，这里的代码不能阻塞 必须开携程")
		}).
		Provide(fxTags.CreateInterFace(memory.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_memory"), nil)).
		Provide(fxTags.CreateInterFace(redis.NewCache, new(cache.Cacher), fxTags.CreateTag("cache_redis"), nil)).
		Mount(controller.NewRouterDog, fxTags.ParamList(fxTags.CreateTag("dogFirst"), fxTags.CreateTag("dogSecond"))).
		//基于接口的使用差别只是Create换成了CreateInterFace 同时显示是哪个接口new(cache.Cacher)
		Mount(controller.NewRouter3,
			fxTags.ParamList(
				fxTags.CreateTag("cache_memory"), fxTags.CreateTag("cache_redis"))).
		Run()
}
