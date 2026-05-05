package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
	"github.com/wike2019/wike_go/v2/func/rateLimiter"
	"github.com/wike2019/wike_go/v2/fxTags"
)

type router struct {
	ctl.Controller
	GetLimit rateLimiter.GetLimit
}

func NewRouter(GetLimit rateLimiter.GetLimit) *router {
	return &router{GetLimit: GetLimit}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OKWithMsg("hello world")
}
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.Use(core.RateLimiter(this.GetLimit, 1, 2))
	//第一个参数	第二个参数	第三个参数
	//限流器 可以根据需要配置不同的限流器	限流器生成令牌桶的速率	限流器的桶容量
	r.GET("/healthz", this.healtzh)
}
func (this router) Path() string {
	return "/"
}
func main() {
	//一个最简单的例子
	g := core.God()
	g.Provide(fxTags.Create(rateLimiter.RateLimit, fxTags.CreateTag("limit"), nil))
	g.Provide(config.Config).Mount(NewRouter, fxTags.ParamList(fxTags.CreateTag("limit"))).Run()
}
