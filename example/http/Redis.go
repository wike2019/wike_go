package http

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	Injector "github.com/shenyisyn/goft-ioc"
	"github.com/wike2019/wike_go/src/core/config"
	"log"
	"time"
)

type RedisConfig struct {
}
func NewRedisConfig() *RedisConfig {
	return &RedisConfig{}
}
func(this *RedisConfig) Redis()  *redis.Client{
	Addr:="192.168.3.2:6379"
	if configData := Injector.BeanFactory.Get((*Config.SysConfig)(nil)); configData != nil {
		Addr = configData.(*Config.SysConfig).Server.Redis
	}
	redisClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     Addr,
		Password: "",                     //密码
		DB:       0,              // redis数据库

		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池数量
		MinIdleConns: 10, //好比最小连接数
		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err!=nil{
		log.Fatal(fmt.Errorf("connect error:%s",err))
	}
	log.Println(pong)
	return redisClient
}
