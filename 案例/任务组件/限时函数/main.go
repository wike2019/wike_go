package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Task"
	"time"
)

func main()  {

	fmt.Println(Task.TaskMustBefore(func(in Task.InChan) {
		in<-"我执行了"
	},time.Second*5))

	fmt.Println(Task.TaskMustBefore(func(in Task.InChan) {
		time.Sleep(time.Second*7)
		in<-"我不能执行了,因为超时了"
	},time.Second*5))

	for {
		time.Sleep(1*time.Second)
	}
}

