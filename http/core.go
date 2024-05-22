package coreHttp

import (
	"github.com/gin-gonic/gin"
	casbinInit "github.com/wike2019/wike_go/lib/casbin"
	"go.uber.org/fx"
)

type GCore struct {
	gin              *gin.Engine //gin引擎
	app              *fx.App     //依赖注入
	provides         []interface{}
	supply           []interface{}
	invokes          []interface{}
	Controller       []interface{}
	middleware       map[string][]gin.HandlerFunc
	globalMiddleware []gin.HandlerFunc
	CronFunc         []map[string]func()
	RoleCtl          *casbinInit.RoleCtl
	StopRun          []func() error
}

const IRoutes = "routes"

func God() *GCore {
	//初始化核心对象
	return &GCore{
		gin:              nil,
		Controller:       make([]interface{}, 0),
		middleware:       map[string][]gin.HandlerFunc{},
		provides:         make([]interface{}, 0),
		invokes:          make([]interface{}, 0),
		globalMiddleware: make([]gin.HandlerFunc, 0),
		CronFunc:         make([]map[string]func(), 0),
		StopRun:          make([]func() error, 0),
	}
}
func (this *GCore) Run() {
	//通过依赖注入调用启动函数
	this.app = fx.New(
		Module,
		fx.NopLogger,
		fx.Provide(fx.Annotate(
			this.NewHTTPServer,
			fx.ParamTags(CreateGroup(IRoutes)), //将路由接口组注入进来
		)),
		fx.Supply(this.supply...),      //注册supply
		fx.Provide(this.provides...),   //注册 provides
		fx.Invoke(this.invokes...),     //注册 invokes
		fx.Provide(this.Controller...), //注册 路由
	)
	this.app.Run() //启动app
}
