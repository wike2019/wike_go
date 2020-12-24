package web

import "github.com/gin-gonic/gin"
// 中间件 接口
type Fairing interface {
	OnRequest(*gin.Context) error
	OnResponse(result interface{}) (interface{}, error)
}
