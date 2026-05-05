package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_app/mysql"
	"github.com/wike2019/wike_app/redis"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ctl"
)

type router struct {
	ctl.Controller
	db    *mysql.MysqlDB
	redis *redis.RedisDB
}

func NewRouter(db *mysql.MysqlDB, redis *redis.RedisDB) *router {
	return &router{db: db, redis: redis}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.OKWithMsg("hello world")
}

func (this *router) doSomething(context *gin.Context) {
	c := this.SetContext(context)
	//这里可以操作mysql和redis
	fmt.Println(this.db)
	fmt.Println(this.redis)
	c.OKWithMsg("hello world")
}

func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
	r.GET("/doSomething", this.doSomething)
}

func (this router) Path() string {
	return "/"
}

func main() {
	g := core.God()
	g.Provide(config.Config, mysql.InitMysql, redis.InitRedis).
		MountWithEmpty(NewRouter).
		Run()
}
