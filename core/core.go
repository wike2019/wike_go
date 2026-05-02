package core

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	constdata "github.com/wike2019/wike_go/v2/constData"
	"github.com/wike2019/wike_go/v2/db"
	"github.com/wike2019/wike_go/v2/fxTags"
	"github.com/wike2019/wike_go/v2/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Controller interface {
	Build(r *gin.RouterGroup, GCore *GCore)
	Path() string
}

type GCore struct {
	gin              *gin.Engine //gin引擎
	app              *fx.App     //依赖注入
	provides         []interface{}
	supply           []interface{}
	invokes          []interface{}
	Controller       []interface{}
	globalMiddleware []gin.HandlerFunc
	CronFunc         []map[string]model.CronJob
	StopRun          []func() error
	Reject           bool
	cfg              *viper.Viper
	db               *db.CoreDb
	Zap              *zap.Logger
}

func God() *GCore {
	//初始化核心对象
	return &GCore{
		gin:              nil,
		Controller:       make([]interface{}, 0),
		provides:         make([]interface{}, 0),
		invokes:          make([]interface{}, 0),
		globalMiddleware: make([]gin.HandlerFunc, 0),
		CronFunc:         make([]map[string]model.CronJob, 0),
		StopRun:          make([]func() error, 0),
		cfg:              nil,
	}
}
func (this *GCore) Run() {

	//通过依赖注入调用启动函数
	this.app = fx.New(
		fx.NopLogger,
		fx.Provide(fx.Annotate(
			this.NewHTTPServer,
			fx.ParamTags(fxTags.CreateGroup(constdata.IRoutes)), //将路由接口组注入进来
		)),
		fx.Supply(this.supply...),      //注册supply
		fx.Provide(this.provides...),   //注册 provides
		fx.Invoke(this.invokes...),     //注册 invokes
		fx.Provide(this.Controller...), //注册 路由
		fx.Supply(this),
	)
	this.app.Run() //启动app
}
