

package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/wike2019/wike_go/http"
	"github.com/wike2019/wike_go/src/core/etcd"
	"github.com/wike2019/wike_go/src/core/ioc"
	"github.com/wike2019/wike_go/src/core/task"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
	"github.com/wike2019/wike_go/src/web"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {

	http.CreateCA()

	signalChan := make(chan os.Signal, 1)
	ioc.New().Config(http.NewDBConfig(),http.NewRedisConfig(),http.NewEtcdConfig())
	app:= web.Ignite(). //初始化脚手架
		Attach(http.NewTokenCheck(), http.NewAddVersion()).
		Mount("", http.NewIndexController())

	ioc.New().ApplyAll()
	task.New().Do("0/2 * * * * *","IndexController.Time()")
	//task.GetCronTask().Start()
	task.Task(func(params ...interface{}) {
		fmt.Println(params)
		time.Sleep(time.Second*5)
	}, func() {
		fmt.Println("任务结束")
	},[]interface{}{10,"wike"})
    // 一定要在 ioc.New().ApplyAll()方法之后不然依赖注入不成功
	//注册服务
	catch:=etcd.EtcdCache()
	id1:= uuid.New().String()
	id2:= uuid.New().String()
	catch.RegService(etcd.ServiceInfo{
		ServiceID: id1 ,
		ServiceName: "wike3",
		ServiceAddr: "http:127.0.0.1:8180",
		ServiceType:  "http",
		ServiceWight: 10,
		ServiceHost:"wike.com",
		Status:LoadBalance.Ready,
	})
	catch.RegService(etcd.ServiceInfo{
		ServiceID:   id2,
		ServiceName: "wike3",
		ServiceAddr: "http:127.0.0.1:8183",
		ServiceType:  "http",
		ServiceWight: 10,
		ServiceHost:"test.com",
		Status:LoadBalance.Ready,
	})


	//wg:=sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	time.Sleep(time.Second*1)
	//	defer wg.Done()
	//	err:=catch.Lock("lock", func(params ...interface{}) {
	//
	//		fmt.Println(params)
	//		time.Sleep(5*time.Second)
	//	},
	//		"wike",1)
	//	fmt.Println(err)
	//}()
	//wg.Add(1)
	//go func() {
	//
	//	defer wg.Done()
	//	err:=catch.Lock("lock", func(params ...interface{}) {
	//
	//		fmt.Println(params)
	//		time.Sleep(5*time.Second)
	//	},
	//		"wike",1)
	//	fmt.Println(err)
	//
	//}()
	//
	//wg.Wait()
	etcd.ReleaseEtcdCache(catch)


	go func() {
	   app.Launch()
   }()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//关闭工作
	<-signalChan
	catch2:=etcd.EtcdCache()
	catch2.UnregService(id1)
	catch2.UnregService(id2)
	etcd.ReleaseEtcdCache(catch2)
}