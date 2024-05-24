package core

import (
	"compress/gzip"
	"context"
	"fmt"
	ginGzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	casbinInit "github.com/wike2019/wike_go/pkg/service/casbin"
	cronInit "github.com/wike2019/wike_go/pkg/service/cron"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func (this *GCore) NewHTTPServer(ControllerList []Controller, lc fx.Lifecycle, zap *zap.Logger, cfg *viper.Viper, DefaultCron *cronInit.DefaultCron, RoleCtl *casbinInit.RoleCtl) *http.Server {
	this.RoleCtl = RoleCtl
	this.zap = zap
	this.cfg = cfg
	r := gin.New()
	this.gin = r                           //缓存gin
	this.gin.MaxMultipartMemory = 32 << 20 // 32 MiB //单次上传总文件最大大小
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
			for _, item := range this.CronFunc {
				for k, v := range item {
					DefaultCron.AddFunc(k, v)
				}
			}
			DefaultCron.DefaultTask()
			go func() {
				DefaultCron.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//清理资源
			defer zap.Sync()
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			this.Reject = true
			DefaultCron.Stop()
			srv.Shutdown(ctx)
			for _, job := range this.StopRun {
				err := job()
				if err != nil {
					zap.Error(fmt.Sprintf("stop func is error: %s\n", err))
				}
			}
			return nil
		},
	})
	return srv
}

func (this *GCore) Default() *GCore {
	this.GlobalUse(ginGzip.Gzip(gzip.DefaultCompression)) //开启压缩
	this.GlobalUse(Reject(this))                          //优雅关闭
	this.GlobalUse(AddTrace())
	this.GlobalUse(CustomRecover(this)) //添加recover中间件和traceId中间件
	this.GlobalUse(AccessLog(this))     //添加c和 跨域中间件
	this.GlobalUse(LimitBodySize(32 << 20))
	return this //添加body数据长度限制中间件
	//全局中间件 //注册用户自定义全局中间件
	//this.gin.Use(rateLimiter.RateLimiter(rateLimiterCache, cfg)) //设置接口根据ip限流
}

// 权限注册函数
func (this *GCore) GetWithRbac(r *gin.RouterGroup, role string, path string, handler gin.HandlerFunc) {
	this.RoleCtl.AddRule(role, r.BasePath(), path, http.MethodGet)
	r.GET(path, handler)
}
func (this *GCore) PostWithRbac(r *gin.RouterGroup, role string, path string, handler gin.HandlerFunc) {
	this.RoleCtl.AddRule(role, r.BasePath(), path, http.MethodPost)
	r.POST(path, handler)
}
func (this *GCore) DelWithRbac(r *gin.RouterGroup, role string, path string, handler gin.HandlerFunc) {
	this.RoleCtl.AddRule(role, r.BasePath(), path, http.MethodDelete)
	r.DELETE(path, handler)
}
func (this *GCore) PutWithRbac(r *gin.RouterGroup, role string, path string, handler gin.HandlerFunc) {
	this.RoleCtl.AddRule(role, r.BasePath(), path, http.MethodPut)
	r.PUT(path, handler)
}
