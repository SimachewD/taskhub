package cache

import (
	"context"
)

func (r *RedisClient) Enqueue(ctx context.Context, queue string, payload string) error {

	return r.client.LPush(ctx, queue, payload).Err()
}

func (r *RedisClient) Dequeue(ctx context.Context,queue string) (string, error) {

	return r.client.RPop(ctx, queue).Result()
}