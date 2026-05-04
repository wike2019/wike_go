package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/v2/cronJob"
	"github.com/wike2019/wike_go/v2/db"
	"github.com/wike2019/wike_go/v2/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func (this *GCore) NewHTTPServer(ControllerList []Controller, db *db.CoreDb, lc fx.Lifecycle, zap *zap.Logger, cfg *viper.Viper, defaultCron *cronJob.DefaultCron) *http.Server {

	this.cfg = cfg
	this.db = db
	this.Zap = zap
	r := gin.New()
	this.gin = r //缓存gin
	maxMultipartMemoryMB := cfg.GetInt64("max_multipart_memory")
	if maxMultipartMemoryMB <= 0 {
		maxMultipartMemoryMB = 20
	}
	this.gin.MaxMultipartMemory = maxMultipartMemoryMB << 20 // 单次上传总文件最大大小
	this.gin.Use(this.globalMiddleware...)
	//健康检查路由
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.String(200, "ok")
	})
	//获取唯一token路由
	r.GET("/api/v1/token", func(c *gin.Context) {
		token := uuid.NewString()
		c.String(200, token)
	})
	for _, route := range ControllerList {
		//注册路由
		group := r.Group(route.Path())
		route.Build(group, this)
	}
	srv := &http.Server{
		Addr:    ":" + cfg.GetString("port"),
		Handler: this.gin,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				zap.Error(err.Error())
				return err
			}
			zap.Debug(fmt.Sprintf("Starting HTTP server at %s", srv.Addr))
			go func() {
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					zap.Error(fmt.Sprintf("HTTP server listen: %s\n", err))
				}
			}()
			for idx, item := range this.CronFunc {
				for k, v := range item {
					cronIdx := idx
					cronKey := k
					fn := v.Fn
					defaultCron.AddFunc(k, func() {
						if this.CronFunc[cronIdx][cronKey].Enabled {
							fn()
						}
					})
					JobTask := &model.CronJob{
						Name:    v.Name,
						Enabled: v.Enabled,
						Cron:    k,
						Func:    runtime.FuncForPC(reflect.ValueOf(v.Fn).Pointer()).Name(),
					}
					this.db.DB.Create(JobTask)
				}

			}

			go func() {
				defaultCron.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//清理资源
			defer zap.Sync()
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			this.Reject = true
			defaultCron.Stop()
			srv.Shutdown(ctx)
			for _, stopJob := range this.StopRun {
				err := stopJob()
				if err != nil {
					zap.Error(fmt.Sprintf("stop func is error: %s\n", err))
				}
			}
			return nil
		},
	})
	return srv
}
