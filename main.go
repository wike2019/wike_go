package main

import (
	"errors"
	"fmt"
	"github.com/wike2019/wike_go/pkg/modules/ants_service"
	"time"
)

//
//type router struct {
//	controller.Controller
//}
//type Invoke struct {
//}
//
//func NewInvoke() *Invoke {
//	return &Invoke{}
//}
//func (this *Invoke) job() error {
//	fmt.Println(time.Now().String())
//	return nil
//}
//func NewRouter() *router {
//	return &router{}
//}
//func (this *router) healtzh(context *gin.Context) {
//	c := this.SetContext(context)
//	re := retry.NewRetry(this.job)
//	re.SetTimes(5)
//	re.SetDelay(time.Second * 1)
//	err := re.Do()
//	c.Success("修改成功", err)
//}
//func (this *router) log(context *gin.Context) {
//	c := this.SetContext(context)
//	c.Success("修改成功", memorylog.LogInfo.All())
//}
//func (this *router) job() error {
//	fmt.Println(time.Now().String())
//	return nil
//}
//func (this *router) stop() error {
//	fmt.Println("这里做局部清理")
//	return nil
//}
//func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
//	GCore.Stop(this.stop)
//	r.GET("healthz", this.healtzh)
//	r.GET("log", this.log)
//}
//func (this router) Path() string {
//	return "/"
//}

type Job struct {
	i2 int
}

func (this *Job) Job() error {
	if this.i2 == 3 {
		return errors.New("错误")
	}
	time.Sleep(time.Duration(this.i2) * time.Second)
	fmt.Println("time: ", this.i2, time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
func main() {

	pool, _ := ants_service.NewPool(2)
	defer pool.Release()
	pool.SetTotal(6)
	for i := 0; i < 6; i++ {
		func(k int) {
			j := &Job{i2: k}
			err := pool.Submit(j.Job)
			if err != nil {
				fmt.Println(err)
			}
		}(i)

	}
	//time.Sleep(1 * time.Second)

	fmt.Println("finish-1")
	pool.Wait() // 等待所有任务完成
	fmt.Println(pool.Ok, pool.Fail, pool.Error())
	fmt.Println("finish")
	return

	//一个最简单的例子
	//g := core.God()
	//g.Stop(func() error {
	//	fmt.Println("这里做全局清理")
	//	return nil
	//})
	//
	//g.Default().GlobalUse(core.CORSMiddleware()) //选择redis作为缓存服务的存储
	//g.Provide(MyViper).Provide(NewInvoke).MountWithEmpty(NewRouter).Invokes(func(r *Invoke) {
	//	go r.job() //这里不能阻塞 所以最好用 go xxx
	//}).Run()
}
