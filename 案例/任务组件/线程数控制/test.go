package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Ioc"
)

type Test struct {

}
//实现bean接口
func (this *Test ) Name() string{
	return  "Test"  //必须和结构体名一样，且不能重复，主要用于以后表达式调用
}

func init() {
	Ioc.New().Beans(New()) //将自己注册到Ioc容器中
}



func New()  *Test {
	return &Test {}
}

func(this *Test ) Show(time int64) {
	fmt.Println(time)
}
