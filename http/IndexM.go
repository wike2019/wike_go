package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/src/web"
)

type TokenCheck struct {}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}
func(this *TokenCheck) OnRequest(ctx *gin.Context) error{
	if ctx.Query("token")==""{
		web.Throw("token requred",503,ctx)
	}
	return nil
}
func(this *TokenCheck) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}

type AddVersion struct {

}
func NewAddVersion() *AddVersion {
	return &AddVersion{}
}
func(this *AddVersion) OnRequest(ctx *gin.Context) error{
	return nil
}
func(this *AddVersion) OnResponse(result interface{}) (interface{}, error){
	if m,ok:=result.(gin.H);ok{
		m["version"]="0.3.0"
		return m,nil
	}
	return result,nil
}

type TestVersion struct {

}
func NewTest() *TestVersion {
	return &TestVersion{}
}
func(this *TestVersion) OnRequest(ctx *gin.Context) error{
	if ctx.Query("wike")==""{
		web.Throw("wike requred",503,ctx)
	}
	return nil
}
func(this *TestVersion) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}

type TestIn struct {

}
func NewTestIn() *TestIn {
	return &TestIn{}
}
func(this *TestIn) OnRequest(ctx *gin.Context) error{
	u:=&User{}
	err:=ctx.BindQuery(u)
	fmt.Println(err)
	if err !=nil{
		web.Throw(err.Error(),503,ctx)
	}
	ctx.Set("_req",u)
	return nil
}
func(this *TestIn) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}