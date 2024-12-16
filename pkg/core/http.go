package core

import (
	"compress/gzip"
	"context"
	"fmt"
	ginGzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	cronInit "github.com/wike2019/wike_go/pkg/service/cron"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"text/template"
	"time"
)

func (this *GCore) NewHTTPServer(ControllerList []Controller, db *CoreDb, lc fx.Lifecycle, zap *zap.Logger, cfg *viper.Viper, defaultCron *cronInit.DefaultCron, roleCtl *RoleCtl) *http.Server {
	this.RoleCtl = roleCtl
	this.zap = zap
	this.db = db
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
					defaultCron.AddFunc(k, v)
				}
			}
			defaultCron.DefaultTask()
			go func() {
				defaultCron.Start()
			}()

			tmpl, err := template.New("markdown").Parse(mdTemplate)
			if err != nil {
				zap.Fatal(err.Error())
			}

			// 创建 Markdown 文件
			file, err := os.Create("接口文档.md")
			if err != nil {
				zap.Fatal(err.Error())
			}

			// 渲染模板到文件
			err = tmpl.Execute(file, this.db.GetData())
			if err != nil {
				zap.Fatal(err.Error())
			}

			err = file.Close()
			if err != nil {
				zap.Fatal(err.Error())
			}
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
