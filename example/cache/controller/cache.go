package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/cache"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

// 注册第三个控制器
type router3 struct {
	ctl.Controller
	redisCache  cache.Cacher
	memoryCache cache.Cacher
	cache       *cache.Service
}

// 注册第二个控制器
func NewRouter3(memoryCache cache.Cacher, redisCache cache.Cacher, cache *cache.Service) *router3 {
	return &router3{redisCache: redisCache, memoryCache: memoryCache, cache: cache}
}
func (this router3) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.POST("/cache", this.cacheData)
}

func (this *router3) cacheData(context *gin.Context) {
	c := this.SetContext(context)
	key := "healtzh"

	data := cache.FindWithCallBack[gin.H](key, time.Second*60, this.cache, func() gin.H {
		return gin.H{"data": "form cache"}
	})
	c.OK(data)
}
func (this router3) Path() string {
	return "/"
}
