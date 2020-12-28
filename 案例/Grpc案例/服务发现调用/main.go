package main

import (
	"context"
	"demo/services"
	"fmt"
	"github.com/wike2019/wike_go/src/Grpc"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/util/LoadBalance"

)

func main()  {

	Ioc.New().Config(NewEtcdConfig())
	Ioc.New().ApplyAll()
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.WithEtcd("wike3",LoadBalance.RoundRobinByWeight,"192.168.3.3"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	rsp,_:=svc.GetUserScore(context.Background(),&services.UserScoreRequest{Users: users})
	fmt.Println(rsp)
}
