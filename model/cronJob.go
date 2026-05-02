package model

type CronJob struct {
	Enabled bool
	Name    string
	Fn      func()
}
