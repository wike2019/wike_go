package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type TokenCheck struct {}
func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}
func(this *TokenCheck) OnRequest(ctx *gin.Context) error{
	fmt.Println("在控制台输出一句话")
	return nil
}
func(this *TokenCheck) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}