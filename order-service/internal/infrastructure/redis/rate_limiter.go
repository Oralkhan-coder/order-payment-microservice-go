package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client      *goredis.Client
	maxRequests int64
	window      time.Duration
}

func NewRateLimiter(client *goredis.Client, maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client:      client,
		maxRequests: int64(maxRequests),
		window:      window,
	}
}

// Allow returns true if the request is within the rate limit, false if exceeded.
func (r *RateLimiter) Allow(ctx context.Context, clientID string) (bool, error) {
	bucket := time.Now().Unix() / int64(r.window.Seconds())
	key := fmt.Sprintf("rate_limit:%s:%d", clientID, bucket)

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return true, err // fail open on Redis error
	}
	if count == 1 {
		r.client.Expire(ctx, key, r.window)
	}

	return count <= r.maxRequests, nil
}
