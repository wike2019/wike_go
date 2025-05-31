package core

import (
	"github.com/wike2019/wike_go/pkg/func/ctl"
)

func (this *GCore) DefaultTask() {
	this.Cron("0 0 3 * * *", func() {
		// 在这里写入具体的任务逻辑
		ctl.ClearChan <- struct{}{}
	}, "每天清除数据任务")
}
