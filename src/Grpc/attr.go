package Grpc

import "google.golang.org/grpc"

type GrpcAttrFunc func(grpc *Grpc) //设置User属性的 函数类型
type GrpcAttrFuncs []GrpcAttrFunc
func(this GrpcAttrFuncs) apply(u *Grpc)  {
	for _,f:=range this{
		f(u)
	}
}
func StreamServerInterceptor(fn grpc.UnaryServerInterceptor) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.ChainServer=append(u.ChainServer,fn)
	}
}
func ServerInterceptor(fn grpc.StreamServerInterceptor) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.ChainStreamServer=append(u.ChainStreamServer,fn)
	}
}
func Ip(ip string) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.Ip=ip
	}
}
func Host(host string) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.Host=host
	}
}
func KeyPath(path string) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.KeyPath=path
	}
}

func WithEtcd(name string,selector int,clientIp string) GrpcAttrFunc  {
	return func(u *Grpc) {
		u.ServerName=name
		u.Selector=selector
		u.ClientIp=clientIp
	}
}