package coreHttp

import (
	"compress/gzip"
	"context"
	"fmt"
	ginGzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/lib/bloom"
	casbinInit "github.com/wike2019/wike_go/lib/casbin"
	"github.com/wike2019/wike_go/lib/rateLimiter"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func (this *GCore) NewHTTPServer(ControllerList []Controller, lc fx.Lifecycle, zap *zap.Logger, cfg *viper.Viper, rateLimiterCache *rateLimiter.RateLimiterCache, Default *cron.Cron, RoleCtl *casbinInit.RoleCtl) *http.Server {
	this.RoleCtl = RoleCtl
	r := gin.New()
	this.gin = r                                          //缓存gin
	r.Use(ginGzip.Gzip(gzip.DefaultCompression))          //开启压缩
	r.Use(AddTrace(), CustomRecover(zap))                 //添加recover中间件和traceId中间件
	r.Use(AccessLog(zap), CORSMiddleware())               //添加日志中间件 和 跨域中间件
	r.Use(LimitBodySize(32 << 20))                        //添加body数据长度限制中间件
	r.MaxMultipartMemory = 32 << 20                       // 32 MiB //单次上传总文件最大大小
	r.Use(this.globalMiddleware...)                       //全局中间件 //注册用户自定义全局中间件
	r.Use(rateLimiter.RateLimiter(rateLimiterCache, cfg)) //设置接口根据ip限流
	//健康检查路由
	r.GET("/api/v1/healthz", func(c *gin.Context) {
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
					Default.AddFunc(k, v)
				}
			}
			go func() {
				Default.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//清理资源
			close(bloom.Clear)
			Default.Stop()
			return srv.Shutdown(ctx)
		},
	})
	return srv
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
