package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)

type JsonController struct {

}



func NewJsonController() *JsonController {
	return &JsonController{}
}
func(this *JsonController) Index (ctx *gin.Context) Web.Json {
	return &Person{ID:"1111",Name:"wike"}
}

func(this *JsonController) Name () string   {
	return "JsonController"
}

func(this *JsonController) Json(ctx *gin.Context) Web.Json {
	return gin.H{"resut":"test"}
}
func(this *JsonController) Build(goft *Web.Goft){
	goft.Handle("GET","/",this.Index)
	goft.Handle("GET","/Json",this.Json)
}