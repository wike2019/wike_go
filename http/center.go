package http

import (
	"github.com/wike2019/wike_go/src/core/redis"
	"log"
)


func NewUserGetter(id string) redis.DBGetterFunc {
	return func() interface{} {
		log.Println("get from db")
		newsModel:=NewUser()
		newsModel.Id=id
		return newsModel
	}
}
