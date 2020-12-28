package main

import (
	"context"
	"demo/services"
	"fmt"
	"github.com/wike2019/wike_go/src/Grpc"
	"io"
	"log"
)

func main()  {
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	for j:=1;j<=3;j++{
		users:=make([]*services.UserInfo,0)
		users=append(users,&services.UserInfo{UserId: 1})
		users=append(users,&services.UserInfo{UserId: 2})
		stream,_:=svc.GetUserScoreByTWS(context.Background())
		stream.Send(&services.UserScoreRequest{Users: users})
		res,err:=stream.Recv()
		if err==io.EOF{
			break
		}
		if err!=nil{
			log.Println(err)
		}
		fmt.Println(res.Users)
	}
}
