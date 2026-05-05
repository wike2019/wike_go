package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisDB struct {
	*redis.Client
}

func InitRedis(cfg *viper.Viper) *RedisDB {
	addr := cfg.GetString("redis")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to ping redis: %s", err.Error())
	}
	log.Println("redis connected:", addr)
	return &RedisDB{Client: client}
}
