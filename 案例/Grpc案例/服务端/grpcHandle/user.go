package grpcHandle

import (
	"context"
	. "demo/services"
	"io"
	"log"
)
type UserService struct {
}
//GetUserScore(context.Context, *UserScoreRequest) (*UserScoreResponse, error)
//GetUserScoreByServerStream(*UserScoreRequest, UserService_GetUserScoreByServerStreamServer) error
//GetUserScoreByClientStream(UserService_GetUserScoreByClientStreamServer) error
//GetUserScoreByTWS(UserService_GetUserScoreByTWSServer) error

//常规调用
func(*UserService) GetUserScore(ctx context.Context,in *UserScoreRequest) (*UserScoreResponse, error) {
	users:=make([]*UserInfo,0)
	for _,user:=range in.Users{
		user.UserScore=10
		users=append(users,user)
	}
	return &UserScoreResponse{Users: users}, nil

}

//服务端流
func(*UserService) 	GetUserScoreByServerStream(in *UserScoreRequest,stream UserService_GetUserScoreByServerStreamServer) error {
	users:=make([]*UserInfo,0)
	for _,user:=range in.Users{
		//做某些操作
		user.UserScore=10
		users=append(users,user)

		//结果发送给客户端
		stream.Send(&UserScoreResponse{Users: users})
	}
	//判断最后是否还有残余数据
	if len(users)>0{
		stream.Send(&UserScoreResponse{Users: users})
	}
	return nil
}

//客户端流
func(*UserService)  GetUserScoreByClientStream(stream UserService_GetUserScoreByClientStreamServer) error{
	users:=make([]*UserInfo,0)
	//循环监听是否有数据
	for{
		//接收数据
		req,err:=stream.Recv()

		//判断函数结束
		if err==io.EOF{ //接收完了
			//把结果给客户端
			return stream.SendAndClose(&UserScoreResponse{Users: users})
		}
		if err!=nil{
			return err
		}
		//处理逻辑
		for _,user:=range req.Users{
			user.UserScore=15  //这里好比是服务端做的业务处理
			users=append(users,user)
		}
	}
}

//双向流
func(*UserService) GetUserScoreByTWS(stream UserService_GetUserScoreByTWSServer) error  {

	users:=make([]*UserInfo,0)
	for{
		//参考  GetUserScoreByClientStream
		req,err:=stream.Recv()
		if err==io.EOF{ //接收完了
			return nil
		}
		if err!=nil{
			return err
		}
		for _,user:=range req.Users{
			user.UserScore=20  //这里好比是服务端做的业务处理
			users=append(users,user)
		}
		err=stream.Send(&UserScoreResponse{Users: users})
		if err!=nil{
			log.Println(err)
		}
		//这里清理数据
		users=(users)[0:0]
	}
}