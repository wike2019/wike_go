package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)

type StringController struct {

}

func NewStringController() *StringController {
	return &StringController{}
}
func(this *StringController) Index(ctx *gin.Context) string   {
	return "111"
}

func(this *StringController) Name () string   {
	return "StringController"
}
func(this *StringController) Build(goft *Web.Goft){
	goft.Handle("GET","/",this.Index)
}