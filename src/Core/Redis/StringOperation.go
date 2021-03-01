package Redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/Core/Bean"

	"time"
)

//专门处理string类型的操作
type RedisStringOperation struct{
	Ctx   context.Context
	Redis *redis.Client `inject:"-"`
}

func init()  {
	Bean.New().Beans(NewStringOperation())
}
func NewStringOperation() *RedisStringOperation {
	return &RedisStringOperation{Ctx: context.Background()}
}
func(this *RedisStringOperation) Name() string{

	return "RedisStringOperation"
}
func(this *RedisStringOperation) Set(key string,value interface{},
					attrs ...*OperationAttr) *Result.ErrorResult {
	exp:=OperationAttrs(attrs).
			 Find(ATTR_EXPIRE).
		     Unwrap_Or([]Result.Any{time.Second*0})[0].(time.Duration)

	nx:=OperationAttrs(attrs).Find(ATTR_NX).Unwrap_Or(nil)[0]
	if nx!=nil{
		return Result.Result(this.Redis.SetNX(this.Ctx,key,value,exp).Result())
	}
	xx:=OperationAttrs(attrs).Find(ATTR_XX).Unwrap_Or(nil)[0]
	if xx!=nil{
		return Result.Result(this.Redis.SetXX(this.Ctx,key,value,exp).Result())
	}
    return Result.Result(this.Redis.Set(this.Ctx,key,value,
   	   exp).Result())

}
func(this *RedisStringOperation) Get(key string ) *Result.ErrorResult {
	 return Result.Result(this.Redis.Get(this.Ctx,key).Result())
}

func(this *RedisStringOperation) MGet(keys ...string )*Result.ErrorResult {
	return Result.Result(this.Redis.MGet(this.Ctx,keys...).Result())
}

