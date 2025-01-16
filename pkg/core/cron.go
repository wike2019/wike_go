package core

import (
	"fmt"
	controller "github.com/wike2019/wike_go/pkg/service/ctl"
	"time"
)

func (this *GCore) DefaultTask() {
	this.Cron("0 0 3 * * *", func() {
		fmt.Println("任务开始执行：", time.Now())
		// 在这里写入具体的任务逻辑
		controller.ClearChan <- struct{}{}
	}, "每天清除数据任务")
}
