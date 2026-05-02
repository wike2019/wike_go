package model

type CronJob struct {
	Enabled bool
	Name    string
	Cron    string
	Func    string
	Fn      func() `gorm:"-"`
}
