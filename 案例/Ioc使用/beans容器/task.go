package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Ioc"
)

type Task struct {

}
//实现bean接口
func (this *Task ) Name() string{
	return  "Task"  //必须和结构体名一样，且不能重复，主要用于以后表达式调用
}

func init() {
	Ioc.New().Beans(New()) //将自己注册到Ioc容器中
}



func New()  *Task {
	return &Task {}
}

func(this *Task ) Show(time int64) {
	fmt.Println(time)
}
