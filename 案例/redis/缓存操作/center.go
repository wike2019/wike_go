package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Redis"
)


func NewUserGetter(id string) Redis.DBGetterFunc {
	return func() interface{} {
		fmt.Println("get from db")
		//各种db操作
		return "数据"+id
	}
}


