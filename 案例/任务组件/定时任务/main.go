package main

import (
	"github.com/wike2019/wike_go/src/core/Task"
	"time"
)

func main()  {
	Test:=&Test{}
	Test.Show(10)//执行一次
	Task.New().Do("0/30 * * * * *","Test.Show(20)") //定时任务
	Task.GetCronTask().Start()




	for {
		time.Sleep(1*time.Second)
	}
}

