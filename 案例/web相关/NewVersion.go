package main

import "github.com/gin-gonic/gin"

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