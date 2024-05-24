package cronInit

import (
	"fmt"
	"github.com/robfig/cron"
	controller "github.com/wike2019/wike_go/pkg/service/http"
	"time"
)

type DefaultCron struct {
	*cron.Cron
}

// 定时任务函数
func NewDefaultCron() *DefaultCron {
	return &DefaultCron{
		Cron: cron.New(),
	}
}
func (this *DefaultCron) DefaultTask() {
	this.Cron.AddFunc("0 0 3 * * *", func() {
		fmt.Println("任务开始执行：", time.Now())
		// 在这里写入具体的任务逻辑
		controller.ClearChan <- struct{}{}
	})
}
