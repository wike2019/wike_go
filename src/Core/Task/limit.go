package Task
type poolChan   chan struct{}
type  TaskLmit struct {
	pool  poolChan
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


