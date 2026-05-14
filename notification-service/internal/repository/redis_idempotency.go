package repository

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisIdempotencyStore struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewRedisIdempotencyStore(client *goredis.Client, ttl time.Duration) *RedisIdempotencyStore {
	return &RedisIdempotencyStore{client: client, ttl: ttl}
}

func (s *RedisIdempotencyStore) MarkIfNew(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return true, nil
	}
	ok, err := s.client.SetNX(ctx, "notif:"+id, "processed", s.ttl).Result()
	return ok, err
}
