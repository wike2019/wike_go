package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
	ctl.Controller
}

func NewRouter() *router {
	return &router{}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OKWithMsg("hello world")
}
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
}
func (this router) Path() string {
	return "/"
}
func main() {
	//一个最简单的例子
	g := core.God()
	// traceId中间件
	g.GlobalUse(core.AddTrace())
	//recover中间件
	g.GlobalUse(core.CustomRecover(g))
	//优雅关闭中间件
	g.GlobalUse(core.Reject(g))
	//日志中间件
	g.GlobalUse(core.AccessLog(g))
	//限制body大小中间件
	g.GlobalUse(core.LimitBodySize(32 << 20))
	//跨域中间件
	g.GlobalUse(core.CORSMiddleware())
	//超时中间件
	g.GlobalUse(core.TimeoutMiddleware(time.Second * 10))
	//自定义中间件
	g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
