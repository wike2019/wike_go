package coreHttp

import "github.com/gin-gonic/gin"

type Controller interface {
	Build(r *gin.RouterGroup, GCore *GCore)
	Path() string
}
