package model

import "gorm.io/gorm"

type Job struct {
	gorm.Model // 主键
	JobSearch
	Cron string `gorm:"size:2000"`
	Func string `gorm:"size:2000"` // 入参数
}
type JobSearch struct {
	Name string `gorm:"size:300" search:"like"` // 路由名称
}

func (this *JobSearch) TableName() string {
	return "jobs"
}

func (this *Job) TableName() string {
	return "jobs"
}
