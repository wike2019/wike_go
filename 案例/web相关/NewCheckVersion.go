package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)

type CheckVersion struct {

}
func NewCheckVersion() *CheckVersion {
	return &CheckVersion {}
}
func(this *CheckVersion) OnRequest(ctx *gin.Context) error{
	if ctx.Query("wike")==""{
		Web.Throw("wike requred",503,ctx)
	}
	return nil
}
func(this *CheckVersion) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}
