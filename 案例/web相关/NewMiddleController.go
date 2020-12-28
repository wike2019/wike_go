package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)

type MiddleController struct {

}

func NewMiddleController() *MiddleController {
	return &MiddleController{}
}
func(this *MiddleController) Index(ctx *gin.Context) string   {
	return "路由级中间件"
}

func(this *MiddleController) Name () string   {
	return "MiddleController"
}
func(this *MiddleController) Build(goft *Web.Goft){
	goft.HandleWithFairing("GET","/index",this.Index, NewCheckName(),NewCheckVersion())
}
