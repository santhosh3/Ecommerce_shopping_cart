package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func RedisRateLimit(redisURL string) (*redis.Client, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Check if the Redis server is reachable by sending a ping
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return rdb, nil
}
