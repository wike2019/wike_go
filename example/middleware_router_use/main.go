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
	time.Sleep(time.Second * 2)
	c.OKWithMsg("hello world")
}
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
	//访问 /healthz 的时候会等候2 秒，访问 /timeOut 的时候会超时 这个就是路由级中间件的作用
	r.Use(core.TimeoutMiddleware(time.Second * 1))
	r.GET("/timeOut", this.healtzh)
}
func (this router) Path() string {
	return "/"
}
func main() {
	//一个最简单的例子
	g := core.God()
	g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
