package utils

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{client: client}
}

func (this *RedisQueue) Push(queueName string, message interface{}) error {
	msg, _ := json.Marshal(message)
	return this.client.LPush(context.Background(), queueName, msg).Err()
}

func (this *RedisQueue) Pop(queueName string, obj interface{}) error {
	message, err := this.client.BLPop(context.Background(), 0, queueName).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(message[1]), obj)
	if err != nil {
		return err
	}
	return nil
}
