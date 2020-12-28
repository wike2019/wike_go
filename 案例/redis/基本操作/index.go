package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Redis"
	"time"
)

type RedisCall struct {
	RedisSet *Redis.RedisStringOperation `inject:"-"`
}

func init() {
	Ioc.New().Beans(new(RedisCall))
}

func(this *RedisCall) Name() string  {
	return "RedisCall"
}
func(this *RedisCall) WithTimeout()    {
	this.RedisSet.Set("keysWithTimeout","我是有过期时间的", Redis.WithExpire(time.Second*10)).Unwrap()
}
func(this *RedisCall) NoTimeout()    {
	this.RedisSet.Set("keysNoTimeout","我是没过期时间的").Unwrap() //设置过期时间
}
func(this *RedisCall) WithNX()    {
	this.RedisSet.Set("keysWithNX","我是设置了NX属性").Unwrap() //设置过期时间
}
func(this *RedisCall) WithXX()    {
	this.RedisSet.Set("keysWithXX","我是设置了属性").Unwrap() //设置过期时间
}

func(this *RedisCall) GetData()    {
	fmt.Println(this.RedisSet.Get("keysNoTimeout").Unwrap() )//取值
}
