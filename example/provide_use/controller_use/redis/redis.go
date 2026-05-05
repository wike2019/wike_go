package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis(cfg *viper.Viper) *redis.Client {
	addr := cfg.GetString("redis")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to ping redis: %s", err.Error())
	}
	log.Println("redis connected:", addr)
	return client
}
