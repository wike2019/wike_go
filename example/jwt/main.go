package main

import (
	"fmt"
	"time"

	"github.com/wike2019/wike_go/v2/func/jwt"
)

const key = "12345567"

type Info struct {
	Name string
}

func main() {
	//一个最简单的例子
	coreData := Info{
		Name: "我是wike",
	}
	//生成一个token
	token, err := jwt.Create[Info](coreData, time.Second*10, key, "wike")
	fmt.Println(token, err)
	info, err := jwt.Parse[Info](token, key)
	fmt.Println(info, err)

}
