package redis

import (
	"regexp"
	"time"
)

type CachePolicy interface {
	Before(key string ) //之前执行
	IfNil(key string,v interface{}) //当空值是操作
	SetOperation(opt *RedisStringOperation) //设置处理器 目前支持string
}

//缓存穿透 策略
type CrossPolicy struct {
	KeyRegx string  //检查key的正则
	Expire time.Duration  //可以配置查找失败过期时间
	opt *RedisStringOperation
}

func NewCrossPolicy(keyRegx string,expire time.Duration) *CrossPolicy {
	return &CrossPolicy{KeyRegx: keyRegx,Expire:expire}
}

func (this *CrossPolicy) Before(key string )  {
		if !regexp.MustCompile(this.KeyRegx).MatchString(key){
			panic("error cache key")
		}
}
func(this *CrossPolicy) IfNil(key string,v interface{})  {
	 	this.opt.Set(key,v,WithExpire(this.Expire)).Unwrap()

}
func(this *CrossPolicy) SetOperation(opt *RedisStringOperation){
	this.opt=opt
}