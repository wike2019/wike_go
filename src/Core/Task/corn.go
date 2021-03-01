package Task

import (
	"github.com/robfig/cron/v3"
	"github.com/shenyisyn/goft-expr/src/expr"
	"github.com/wike2019/wike_go/src/Core/Bean"
	"log"
	"sync"
)

var taskCron *cron.Cron
var CronInstance *SetInterVal
type  SetInterVal struct {

}
var setInterVal sync.Once
var onceCron sync.Once
func New()  *SetInterVal{
	setInterVal.Do(func() {
		CronInstance=&SetInterVal{}
	})
	return CronInstance
}

func (this *SetInterVal)Do(cron string, exprorfunc interface{}) *SetInterVal {
	var err error
	if f, ok := exprorfunc.(func()); ok {
		_, err = GetCronTask().AddFunc(cron, f)
	} else if exprorfunc, ok := exprorfunc.(string); ok {
		_, err = GetCronTask().AddFunc(cron, func() {
			expr.BeanExpr(exprorfunc, Bean.New().ExprData)
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}


//初始化定时任务器
func GetCronTask() *cron.Cron {

	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}