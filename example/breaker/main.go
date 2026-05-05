package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/breaker"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
	ctl.Controller
	Breaker *gobreaker.CircuitBreaker
}

func NewRouter() *router {
	return &router{
		Breaker: breaker.NewCircuitBreaker("breaker", 5, time.Second*10, time.Second*3),
	}
}

func (this *router) Job() (interface{}, error) {

	// 模拟一个 HTTP 请求
	resp, err := http.Get("https://baidu.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	return "Request succeeded", nil

}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	data, err := this.Breaker.Execute(this.Job)
	ctl.Error(err, 500)
	c.OK(data)
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
	g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
