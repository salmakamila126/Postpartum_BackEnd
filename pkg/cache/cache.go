package cache

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewFromEnv() (*Cache, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, nil
	}

	db := 0
	if raw := os.Getenv("REDIS_DB"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		db = value
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{client: client}, nil
}

func (c *Cache) Enabled() bool {
	return c != nil && c.client != nil
}

func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	if !c.Enabled() {
		return false, nil
	}

	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !c.Enabled() {
		return nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, raw, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if !c.Enabled() || len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}
