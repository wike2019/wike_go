package core

import (
	"github.com/gin-gonic/gin"
	casbinInit "github.com/wike2019/wike_go/pkg/service/casbin"
	cronInit "github.com/wike2019/wike_go/pkg/service/cron"
	zaplog "github.com/wike2019/wike_go/pkg/service/log"
	"go.uber.org/fx"
	"net/http"
	"reflect"
)

var Module = fx.Module("infra",
	fx.Provide(zaplog.GetLogger), //默认日志
	fx.Provide(cronInit.NewDefaultCron, casbinInit.NewEnforcer, casbinInit.NewCtl), //定时器任务 rbac权限
	fx.Invoke(func(*http.Server) {}),
	fx.Provide(InitDb),
)

// 用于没有参数的依赖注入
func (this *GCore) Config(cfgs ...interface{}) *GCore {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfg)
		if t.Kind() != reflect.Ptr {
			panic("required ptr object") //必须是指针对象
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		} //处理依赖注入 (new)
		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)

			if callRet != nil && len(callRet) == 1 {
				this.supply = append(this.supply, callRet[0].Interface())
			}
		}
	}
	return this
}

// 用于有注入参数的依赖注入
func (this *GCore) Provide(list ...interface{}) *GCore {
	this.provides = append(this.provides, list...)
	return this
}

// 用于主动调用
func (this *GCore) Invokes(list ...interface{}) *GCore {
	this.invokes = append(this.invokes, list...)
	return this
}

// 用于注册全局中间件
func (this *GCore) GlobalUse(middleware ...gin.HandlerFunc) *GCore {
	this.globalMiddleware = append(this.globalMiddleware, middleware...)
	return this
}

// 用于挂载带参数控制器
func (this *GCore) Mount(class interface{}, params []string) *GCore {
	this.Controller = append(this.Controller, CreateInterFace(class, new(Controller), CreateGroup("routes"), params))
	return this
}

// 用于挂载不带参数控制器
func (this *GCore) MountWithEmpty(class interface{}) *GCore {
	this.Controller = append(this.Controller, CreateInterFace(class, new(Controller), CreateGroup("routes"), []string{}))
	return this
}

// 用于挂载已经存在带对象
func (this *GCore) Supply(supply ...interface{}) *GCore {
	this.supply = append(this.supply, supply...)
	return this
}

// 用于添加定时任务
func (this *GCore) Cron(spec string, cmd func()) *GCore {
	this.CronFunc = append(this.CronFunc, map[string]func(){spec: cmd})
	return this
}
func (this *GCore) Stop(job func() error) *GCore {
	this.StopRun = append(this.StopRun, job)
	return this
}
