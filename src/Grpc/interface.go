package Grpc

import (
	"google.golang.org/grpc"
	)

type Middle interface{
	Action() grpc.UnaryServerInterceptor
}
type MiddleStream interface{
	Action() grpc.StreamServerInterceptor
}