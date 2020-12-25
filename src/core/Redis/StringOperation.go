package Redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/core/Ioc"

	"time"
)

//专门处理string类型的操作
type RedisStringOperation struct{
	ctx context.Context
	Redis  *redis.Client `inject:"-"`
}

func init()  {
	Ioc.New().Beans(NewStringOperation())
}
func NewStringOperation() *RedisStringOperation {
	return &RedisStringOperation{ctx:context.Background()}
}
func(this *RedisStringOperation) Name() string{
	return "RedisStringOperation"
}
func(this *RedisStringOperation) Set(key string,value interface{},
					attrs ...*OperationAttr) *Result.ErrorResult {
	exp:=OperationAttrs(attrs).
			 Find(ATTR_EXPIRE).
		     Unwrap_Or(time.Second*0).(time.Duration)

	nx:=OperationAttrs(attrs).Find(ATTR_NX).Unwrap_Or(nil)
	if nx!=nil{
		return Result.Result(this.Redis.SetNX(this.ctx,key,value,exp).Result())
	}
	xx:=OperationAttrs(attrs).Find(ATTR_XX).Unwrap_Or(nil)
	if xx!=nil{
		return Result.Result(this.Redis.SetXX(this.ctx,key,value,exp).Result())
	}
    return Result.Result(this.Redis.Set(this.ctx,key,value,
   	   exp).Result())

}
func(this *RedisStringOperation) Get(key string ) *Result.ErrorResult {
	 return Result.Result(this.Redis.Get(this.ctx,key).Result())
}

func(this *RedisStringOperation) MGet(keys ...string )*Result.ErrorResult {
	return Result.Result(this.Redis.MGet(this.ctx,keys...).Result())
}

