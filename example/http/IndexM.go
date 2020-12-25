package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/wike2019/wike_go/src/Web"
	"github.com/wike2019/wike_go/src/util/Validate"
)

type TokenCheck struct {}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}
func(this *TokenCheck) OnRequest(ctx *gin.Context) error{
	if ctx.Query("token")==""{
		Web.Throw("token requred",503,ctx)
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
		Web.Throw("wike requred",503,ctx)
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
	binding.Form.Bind(ctx.Request,u)
	//验证器的使用
	err:=Validate.New().Validate.Struct(u)
	fmt.Println(err)
	if err !=nil{
		Web.Throw(Validate.New().Msg(u,err),503,ctx)
	}
	ctx.Set("_req",u)
	return nil
}
func(this *TestIn) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}