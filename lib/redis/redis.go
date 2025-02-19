// Package redis provides simple Redis cache functions.
package redis

import (
	"fmt"

	"github.com/bhbdev/jam/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func Setup(cfg *config.Redis) (*RedisClient, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address(),
		DB:   cfg.Db,
	})
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	return &RedisClient{client}, nil

}
