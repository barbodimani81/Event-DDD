package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	Client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{Client: client}
}

func (r *RedisRateLimiter) IsAllowed(ctx context.Context, userID string, limit int) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", userID)

	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error incrementing key: %w", err)
	}

	if count == 1 {
		r.Client.Expire(ctx, key, 60*time.Second)
	}

	if count > int64(limit) {
		return false, fmt.Errorf("rate limit reached")
	}

	return true, nil
}
