package main

import (
	"context"
	"demo/services"
	"fmt"
	"github.com/wike2019/wike_go/src/Grpc"
	"io"
)

func main()  {
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	stream,_:=svc.GetUserScoreByServerStream(context.Background(),&services.UserScoreRequest{Users: users})
	for{
		res,err:=stream.Recv()
		if err==io.EOF{
			break
		}
		fmt.Println(res.Users)

	}
}
