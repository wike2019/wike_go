package cronInit

import (
	"github.com/robfig/cron"
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
