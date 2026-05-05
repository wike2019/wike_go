package controller

import (
	"fmt"

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
}

// 注册第二个控制器
func NewRouter3(memoryCache cache.Cacher, redisCache cache.Cacher) *router3 {
	return &router3{redisCache: redisCache, memoryCache: memoryCache}
}
func (this router3) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.POST("/cache", this.cache)
}

func (this *router3) cache(context *gin.Context) {
	c := this.SetContext(context)
	//这里可以使用不同的cache对象
	this.memoryCache.Set("test", "1", 0)
	str := ""
	this.memoryCache.Get("test", &str)
	fmt.Println(str)

	this.redisCache.Set("test2", "2", 0)
	str2 := ""
	this.redisCache.Get("test2", &str2)
	fmt.Println(str2)
	c.OK("ok")
}
func (this router3) Path() string {
	return "/"
}
