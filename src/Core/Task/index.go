package Task

func init() {
	taskList := getTaskList() //得到任务列表
	go func() {
		for t := range taskList {
			doTask(t)
		}
	}()
}




