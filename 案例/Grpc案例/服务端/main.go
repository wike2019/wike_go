package main

import (
	"context"
	"demo/grpcHandle"
	"demo/services"
	"github.com/google/uuid"
	"github.com/wike2019/wike_go/src/Grpc"
	"github.com/wike2019/wike_go/src/core/Etcd"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	signalChan := make(chan os.Signal, 1)
	Ioc.New().Config(NewEtcdConfig())
	Ioc.New().ApplyAll()
	Instance:= Etcd.EtcdCache()
	id1:= uuid.New().String()
	defer Etcd.ReleaseEtcdCache(Instance)
	//id2:= uuid.New().String()
	Instance.RegService(Etcd.ServiceInfo{
		ServiceID: id1 ,
		ServiceName: "wike3",
		ServiceAddr: "192.168.3.3:8081",
		ServiceType:  "http",
		DefaultBalance:&LoadBalance.DefaultBalance{Weight:10,Status:LoadBalance.Ready,Id:id1},
		ServiceHost:"wike.com",

	})
	rpcServer:=Grpc.NewServer(Grpc.KeyPath("./keys"),
		Grpc.StreamServerInterceptor(func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			log.Printf("调用前")
			err:=handler(srv, stream)
			log.Printf("调用前后")
			return  err
	}),
		Grpc.ServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			log.Printf("调用前")
			resp,err:=handler(ctx, req)
			log.Printf("调用前后")
			return resp, err
		}))
	services.RegisterUserServiceServer(rpcServer,new(grpcHandle.UserService))
	go func() {
		lis,_:=net.Listen("tcp","192.168.3.3:8081")
		rpcServer.Serve(lis)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//关闭工作
	<-signalChan
	Instance2:= Etcd.EtcdCache()
	Instance2.UnregService(id1)
	//catch2.UnregService(id2)
	Etcd.ReleaseEtcdCache(Instance2)

}
