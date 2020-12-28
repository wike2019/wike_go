package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/Web"
)
type CheckName struct {

}
func NewCheckName() *CheckName {
	return &CheckName {}
}
func(this *CheckName ) OnRequest(ctx *gin.Context) error{
	if ctx.Query("wike")==""{
		Web.Throw("wikerequred",503,ctx)
	}
	return nil
}
func(this *CheckName ) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}
