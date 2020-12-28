package main



import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"github.com/wike2019/wike_go/src/core/Redis"
)
//专门处理string类型的操作
type RedisStringOperation struct{
	Ctx context.Context
	Redis  *redis.Client `inject:"-"`
}

func main()  {

	Ioc.New().Config(NewDBConfig(),NewRedisConfig()) //注册db对象和redis对象
	Ioc.New().ApplyAll()//执行注入功能
	opt:=&Redis.RedisStringOperation{}
	opt=Ioc.New().Get(opt).(*Redis.RedisStringOperation)
	fmt.Println(opt.Redis)//此时有值
}




