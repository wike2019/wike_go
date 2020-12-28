package main

import (
	"github.com/wike2019/wike_go/src/core/Task"
	"time"
)

func main()  {
	Test:=&Test{}

	limit:=Task.LimtTask(5)
	i:=0
	for  {
		i++
		go func() {
			limit.CanDo()
			defer limit.Close()
			Test.Show(10)//执行一次
			time.Sleep(time.Second*5)
		}()
		if i>100 {
			break
		}
	}
	
	for {
		time.Sleep(1*time.Second)
	}
}

