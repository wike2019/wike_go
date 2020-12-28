package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wike2019/wike_go/src/Web"
	"github.com/wike2019/wike_go/src/core/sql"
)

type SqlController struct {
	Db *gorm.DB                          `inject:"-"`
}


func NewSqlController() *SqlController {
	return &SqlController{}
}
func(this *SqlController) Index (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id > ?").WithArgs(1)
}
func(this *SqlController) Single (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst()
}
func(this *SqlController) WithKey (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst().WithKey("result")
}

func(this *SqlController) WithChange  (ctx *gin.Context) Web.Json {
	ret:=  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst().Get()
	m := ret.(map[string]interface{})
	m["additive"]="添加一些数据"
	return  m
}
func(this *SqlController) Name () string   {
	return "SqlController"
}
func(this *SqlController) Build(goft *Web.Goft){
	goft.Handle("GET","/",this.Index)
	goft.Handle("GET","/Single",this.Single)
	goft.Handle("GET","/WithKey",this.WithKey)
	goft.Handle("GET","/WithChange",this.WithChange)
}