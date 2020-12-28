package main

import (
	"github.com/wike2019/wike_go/src/core/Ioc"
)

func main()  {
	Ioc.New().Config(NewRedisConfig())
	Ioc.New().ApplyAll()
	call:=new(RedisCall)
	call=Ioc.New().Get(call).(*RedisCall)
	call.WithTimeout()
	call.WithXX()
	call.WithNX()
	call.NoTimeout()
	call.GetData()
}
