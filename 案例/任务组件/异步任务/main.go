package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Task"
	"time"
)

func main()  {


	Task.Task( //延迟任务 用于处理一些后置操作
	func(params ...interface{}) {
		fmt.Println(params)
		time.Sleep(time.Second*5)  //处理函数
	},
	func() {//回调
		fmt.Println("任务结束")
	},
	[]interface{}{10,"wike"},//参数
	)

	
	for {
		time.Sleep(1*time.Second)
	}
}

