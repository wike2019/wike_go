package coreHttp

import (
	"bytes"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wike2019/wike_go/lib/controller"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Reject(god *GCore) gin.HandlerFunc {
	return func(context *gin.Context) {
		if god.Reject {
			context.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "服务暂时关闭", "code": 503})
			return
		}
		context.Next()
	}
}

func AddTrace() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := ""
		traceIdHeader := context.Request.Header.Get("trace_id")
		traceIdQuery := context.Query("trace_id")
		idempotent := context.Query("idempotent")
		if idempotent == "check" && traceIdQuery == "" {
			context.AbortWithStatusJSON(200, gin.H{"message": "非法数据", "code": 200, "trace_id": "0"})
			return
		}
		if traceIdQuery != "" {
			traceId = traceIdQuery
		} else {
			traceId = traceIdHeader
		}
		if traceId == "" {
			traceId = uuid.NewString()
		}
		context.Set("trace_id", traceId)
		context.Next()
	}
}

// 日志中间件
func AccessLog(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		var reqBody []byte
		if ctx.Request.Body != nil {
			reqBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		ctx.Next()
		latency := time.Now().Sub(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		logger.Info("接口访问日志",
			zap.String("path", path),
			zap.String("method", method),
			zap.String("http_host", ctx.Request.Host),
			zap.String("ua", ctx.Request.UserAgent()),
			zap.String("remote_addr", ctx.Request.RemoteAddr),
			zap.Int("request_body_size", len(reqBody)),
			zap.Int("status_code", statusCode),
			zap.Int("error_code", ctx.GetInt("error_code")),
			zap.String("error_msg", ctx.GetString("error_code")),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// 跨域中间
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// 异常恢复中间件
func CustomRecover(logger *zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				check, ok := err.(*controller.StatusError)
				if ok {
					logger.Warn("主动抛出的错误", zap.String("path", c.Request.URL.Path), zap.String("error", check.Msg), zap.Int("code", check.Code))
					c.AbortWithStatusJSON(check.Code, gin.H{"message": check.Msg, "code": check.Code, "trace_id": c.GetString("trace_id")})
					return
				}
				logger.Error("接口错误", zap.String("path", c.Request.URL.Path), zap.String("error", check.Msg), zap.Int("code", check.Code))
				c.AbortWithStatusJSON(500, gin.H{"message": "Internal Server Error", "code": 500, "trace_id": c.GetString("trace_id")})
				return
			}
		}()
		c.Next()
	}
}

// body数据大小限制中间件
func LimitBodySize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentLength := c.Request.Header.Get("Content-Length")
		length, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			c.Next()
			return
		}
		if length > maxSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "请求体内容太长", "trace_id": c.GetString("trace_id")})
			return
		}
		c.Next()
	}
}

// 超时中间件
func TimeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(TimeOutResponse),
	)
}

// 超时响应
func TimeOutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code":     http.StatusGatewayTimeout,
		"msg":      "系统超时了",
		"trace_id": c.GetString("trace_id"),
	})
}

// RBAC中间件
func Authorizer(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		ok, _ := e.Enforce(role, c.FullPath(), c.Request.Method)
		if !ok {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"code":     http.StatusNonAuthoritativeInfo,
				"msg":      "当前角色没有访问权限",
				"trace_id": c.GetString("trace_id"),
			})
			c.Abort()
		}
	}
}
