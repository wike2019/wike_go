package cronInit

import (
	"github.com/robfig/cron"
)

// 定时任务函数
func DefaultCron() *cron.Cron {
	return cron.New()
}
