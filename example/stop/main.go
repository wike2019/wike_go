package main

import (
	"fmt"

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
	g.Provide(config.Config).MountWithEmpty(NewRouter).Stop(func() error {
		fmt.Println("这里做全局清理")
		return nil
	}).Run()
}
