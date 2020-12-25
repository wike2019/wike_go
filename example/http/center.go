package http

import (
	"github.com/wike2019/wike_go/src/core/Redis"
	"log"
)


func NewUserGetter(id string) Redis.DBGetterFunc {
	return func() interface{} {
		log.Println("get from db")
		newsModel:= NewUser()
		newsModel.Id=id
		return newsModel
	}
}

func NewUserGetter2(id string) Redis.DBGetterFunc {
	return func() interface{} {
		log.Println("get from db2")

		return "11111111"
	}
}
