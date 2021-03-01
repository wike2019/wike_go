package Task

import "sync"

var taskList chan *TaskExecutor //任务列表
var taskLmit *TaskLmit  //限制器用于限制延迟任务,防止内存爆了
var once sync.Once

type TaskExecutor struct {
	f        TaskFunc
	p        []interface{} //参数
	callback func()
}

type TaskFunc func(params ...interface{})



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

func getTaskList() chan *TaskExecutor {
	once.Do(func() {
		taskList = make(chan *TaskExecutor) //初始化 做了限制
		taskLmit =LimtTask(10000)  //设置携程数为8000以内
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
