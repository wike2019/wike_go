package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/pkg/core"
	controller "github.com/wike2019/wike_go/pkg/service/http"
	"github.com/wike2019/wike_go/pkg/service/memorylog"
	"github.com/wike2019/wike_go/pkg/service/retry"
	"time"
)

type router struct {
	controller.Controller
}

func NewRouter() *router {
	return &router{}
}
func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	re := retry.NewRetry(this.job)
	re.SetTimes(5)
	re.SetDelay(time.Second * 1)
	err := re.Do()
	c.Success("修改成功", err)
}
func (this *router) log(context *gin.Context) {
	c := this.SetContext(context)
	c.Success("修改成功", memorylog.LogInfo.All())
}
func (this *router) job() error {
	fmt.Println(time.Now().String())
	return nil
}
func (this *router) stop() error {
	fmt.Println("这里做局部清理")
	return nil
}
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	GCore.Stop(this.stop)
	r.GET("healthz", this.healtzh)
	r.GET("log", this.log)
}
func (this router) Path() string {
	return "/"
}

func main() {

	//一个最简单的例子
	g := core.God()
	g.Stop(func() error {
		fmt.Println("这里做全局清理")
		return nil
	})
	g.Default().GlobalUse(core.CORSMiddleware()) //选择redis作为缓存服务的存储
	g.Provide(MyViper).MountWithEmpty(NewRouter).Invokes(func(r *router) {
		go r.job() //这里不能阻塞 所以最好用 go xxx
	}).Run()
}
