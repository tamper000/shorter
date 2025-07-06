package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"urlshort/internal/config"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(config config.Redis) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, err
	}

	cache := &Cache{
		client: rdb,
		ttl:    time.Duration(config.TTL) * time.Minute,
	}

	return cache, nil
}

func (c *Cache) AddCache(alias, link string) error {
	ctx := context.Background()
	err := c.client.Set(ctx, alias, link, c.ttl).Err()
	return err
}

func (c *Cache) GetCache(alias string) (string, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, alias).Result()
	return val, err
}

func (c *Cache) DeleteCache(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	return err
}

func (c *Cache) CloseCache() {
	c.client.Close()
}
