package core

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	constdata "github.com/wike2019/wike_go/v2/constData"
	"github.com/wike2019/wike_go/v2/cronJob"
	"github.com/wike2019/wike_go/v2/db"
	"github.com/wike2019/wike_go/v2/fxTags"
	"github.com/wike2019/wike_go/v2/model"
	"go.uber.org/fx"
)

var Module = fx.Module("infra",
	fx.Provide(NewLogger),              //日志
	fx.Provide(cronJob.NewDefaultCron), //定时器任务
	fx.Invoke(func(*http.Server) {}),
	fx.Provide(db.InitDb),
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
	this.Controller = append(this.Controller, fxTags.CreateInterFace(class, new(Controller), fxTags.CreateGroup(constdata.IRoutes), params))
	return this
}

// 用于挂载不带参数控制器
func (this *GCore) MountWithEmpty(class interface{}) *GCore {
	this.Controller = append(this.Controller, fxTags.CreateInterFace(class, new(Controller), fxTags.CreateGroup(constdata.IRoutes), []string{}))
	return this
}

// 用于挂载已经存在带对象
func (this *GCore) Supply(supply ...interface{}) *GCore {
	this.supply = append(this.supply, supply...)
	return this
}

// 用于添加定时任务
func (this *GCore) Cron(spec string, cmd func(), Job string, enabled bool) *GCore {
	for _, item := range this.CronFunc {
		for _, v := range item {
			if v.Name == Job {
				panic(fmt.Sprintf("cron job name %q already exists", Job))
			}
		}
	}
	this.CronFunc = append(this.CronFunc, map[string]model.CronJob{spec: {Enabled: enabled, Name: Job, Fn: cmd}})

	return this
}
func (this *GCore) Stop(job func() error) *GCore {
	this.StopRun = append(this.StopRun, job)
	return this
}

func (this *GCore) StopTask(name string) *GCore {
	this.cronMu.Lock()
	for index, _ := range this.CronFunc {
		for k, job := range this.CronFunc[index] {
			if job.Name == name {
				job.Enabled = false
				this.CronFunc[index][k] = job
			}
		}
	}
	this.cronMu.Unlock()
	this.db.DB.Model(&model.CronJob{}).Where("name = ?", name).Update("enabled", false)
	return this
}

func (this *GCore) RestartTask(name string) *GCore {
	this.cronMu.Lock()
	for index, _ := range this.CronFunc {
		for k, job := range this.CronFunc[index] {
			if job.Name == name {
				job.Enabled = true
				this.CronFunc[index][k] = job
			}
		}
	}
	this.cronMu.Unlock()
	this.db.DB.Model(&model.CronJob{}).Where("name = ?", name).Update("enabled", true)
	return this
}
