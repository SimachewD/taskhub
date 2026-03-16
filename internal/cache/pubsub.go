package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func (r *RedisClient) Publish(channel string, message string) error {
	return r.client.Publish(context.Background(), channel, message).Err()
}

func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(context.Background(), channel)
}