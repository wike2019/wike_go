package Task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/shenyisyn/goft-expr/src/expr"
	"github.com/wike2019/wike_go/src/core/ioc"
	"log"
	"sync"
	"time"
)

type  SetInterVal struct {

}
type  TaskLmit struct {
	pool  poolChan
}
type poolChan   chan struct{}

type TaskExecutor struct {
	f        TaskFunc
	p        []interface{} //参数
	callback func()
}

type TaskFunc func(params ...interface{})
type InChan chan interface{}

var taskCron *cron.Cron
var taskList chan *TaskExecutor //任务列表
var taskLmit *TaskLmit  //限制器用于限制延迟任务,防止内存爆了
var CronInstance *SetInterVal

var once sync.Once
var setInterVal sync.Once
var onceCron sync.Once


func init() {
	taskList := getTaskList() //得到任务列表
	go func() {

		for t := range taskList {
			doTask(t)
		}
	}()
}
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
			 expr.BeanExpr(exprorfunc, Ioc.New().ExprData)
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}
func doTask(t *TaskExecutor) {
	go func() {
		defer func() {
			taskLmit.Close()
			if t.callback != nil {
				t.callback()
			}
		}()
		t.Exec()
	}()
}
//初始化定时任务器
func GetCronTask() *cron.Cron {

	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}
func getTaskList() chan *TaskExecutor {
	once.Do(func() {
		taskList = make(chan *TaskExecutor) //初始化 做了限制
		taskLmit =LimtTask(8000)  //设置携程数为8000以内
	})
	return taskList
}


//辅助函数 初始化任务器
func NewTaskExecutor(f TaskFunc, p []interface{}, callback func()) *TaskExecutor {
	return &TaskExecutor{f: f, p: p, callback: callback}
}
func (this *TaskExecutor) Exec() { //执行任务
	this.f(this.p...)
}
func Task(f TaskFunc, cb func(), params ...interface{}) {
	if f == nil {
		return
	}
	go func() {
		taskLmit.CanDo()
		getTaskList() <- NewTaskExecutor(f, params, cb)//增加任务队列
	}()
}


//控制任务数
func  LimtTask(maxNum int) *TaskLmit{
	task :=&TaskLmit{}
	task.pool=make(chan struct{},maxNum)
	return task
}


func (this * TaskLmit) CanDo() {
	this.pool<-struct {}{}
}
func (this * TaskLmit) Close() {
	<-this.pool
}


func TaskMustBefore(job func(in InChan),d time.Duration) (interface{},  error) {
	ret:=make(InChan)
	go job(ret)
	select {
	case r:=<-ret:
		return r,nil
	case <-time.After(d):
		return nil,fmt.Errorf("time out")
	}
}