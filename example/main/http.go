//
//
package main



//import (
//	"fmt"
//	"github.com/wike2019/wike_go/example/http"
//	"github.com/wike2019/wike_go/src/Grpc"
//	"github.com/wike2019/wike_go/src/core/Ioc"
//	"github.com/wike2019/wike_go/src/core/Redis"
//)
//
//func main()  {
//	Ioc.New().Config(http.NewDBConfig(),http.NewRedisConfig()) //注册db对象和redis对象
//	Ioc.New().ApplyAll()//执行注入功能
//	opt:=&Redis.RedisStringOperation{}
//	opt=Ioc.New().Get(opt).(*Redis.RedisStringOperation)
//	fmt.Println(opt.Redis)//此时有值
//}
//
import (
	"fmt"
	"github.com/google/uuid"
	grpcHandle "github.com/wike2019/wike_go/example/grpc/handle"
	"github.com/wike2019/wike_go/example/http"
	"github.com/wike2019/wike_go/example/services"
	"github.com/wike2019/wike_go/src/Grpc"
	"github.com/wike2019/wike_go/src/Web"
	"github.com/wike2019/wike_go/src/core/Etcd"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Task"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
	"github.com/wike2019/wike_go/src/util/Validate"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main()  {

	http.CreateCA()
	Validate.New().AddValiDate("NameRequired", http.NameRequired)
	Validate.New().AddValiDate("CheckName", http.CheckName)
	signalChan := make(chan os.Signal, 1)
	Ioc.New().Config(http.NewDBConfig(), http.NewRedisConfig(), http.NewEtcdConfig())
	app:= Web.New(). //初始化脚手架
		Attach(http.NewTokenCheck(), http.NewAddVersion()).
		Mount("", http.NewIndexController())

	Ioc.New().ApplyAll()
	Task.New().Do("0/2 * * * * *","IndexController.Time()")
	//task.GetCronTask().Start()
	Task.Task(func(params ...interface{}) {
		fmt.Println(params)
		time.Sleep(time.Second*5)
	}, func() {
		fmt.Println("任务结束")
	},[]interface{}{10,"wike"})
  //// 一定要在 ioc.New().ApplyAll()方法之后不然依赖注入不成功
	////注册服务
	catch:= Etcd.EtcdCache()
	id1:= uuid.New().String()
	//id2:= uuid.New().String()
	catch.RegService(Etcd.ServiceInfo{
		ServiceID: id1 ,
		ServiceName: "wike3",
		ServiceAddr: "192.168.3.3:8081",
		ServiceType:  "http",
		DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:id1},
		ServiceHost:"wike.com",

	})
	//catch.RegService(Etcd.ServiceInfo{
	//	ServiceID: id2 ,
	//	ServiceName: "wike3",
	//	ServiceAddr: "http:127.0.0.1:81830",
	//	ServiceType:  "http",
	//	DefaultBalance:&LoadBalance.DefaultBalance{Weight:20,Status:LoadBalance.Ready},
	//	ServiceHost:"wike.com",
	//
	//})
  // //管道操作
	//http.Do()
  //
  //
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
	//Etcd.ReleaseEtcdCache(catch)

	rpcServer:=Grpc.NewServer(Grpc.KeyPath("./keys"))
	services.RegisterUserServiceServer(rpcServer,new(grpcHandle.UserService))
	go func() {
	   app.Launch()
  }()
	go func() {
		lis,_:=net.Listen("tcp","192.168.3.3:8081")
		rpcServer.Serve(lis)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//关闭工作
	<-signalChan
	catch2:= Etcd.EtcdCache()
	catch2.UnregService(id1)
	//catch2.UnregService(id2)
	Etcd.ReleaseEtcdCache(catch2)
}