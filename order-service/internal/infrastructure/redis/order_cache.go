package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Oralkhan-coder/order-service/internal/model"
	goredis "github.com/redis/go-redis/v9"
)

type OrderCache struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewOrderCache(client *goredis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{client: client, ttl: ttl}
}

// Get returns the cached order, or (nil, nil) on a cache miss.
func (c *OrderCache) Get(ctx context.Context, id string) (*model.Order, error) {
	data, err := c.client.Get(ctx, cacheKey(id)).Bytes()
	if err == goredis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (c *OrderCache) Set(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, cacheKey(order.ID), data, c.ttl).Err()
}

func (c *OrderCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, cacheKey(id)).Err()
}

func cacheKey(id string) string {
	return "order:" + id
}
