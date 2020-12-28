package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)

type StringController struct {

}

func NewStringController() *ErrorController {
	return &ErrorController{}
}
func(this *ErrorController) Error(ctx *gin.Context) string   {

	Web.Error(fmt.Errorf("抛出一个错误"),"这个是返回的提示消息")
	return "111"
}
func(this *ErrorController) Throw(ctx *gin.Context) string   {

	Web.Throw("错误消息",500,ctx)
	return "111"
}

func(this *ErrorController) Name () string   {
	return "ErrorController"
}
func(this *ErrorController) Build(goft *Web.Goft){
	goft.Handle("GET","/Error",this.Throw)
	goft.Handle("GET","/Throw",this.Throw)
}