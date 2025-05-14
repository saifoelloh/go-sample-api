package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type RedisOptions struct {
	Addr string
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(opt RedisOptions) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{Addr: opt.Addr})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.InternalServerError(
			errors.WithScope("RedisClient"),
			errors.WithLocation("NewRedisClient.Ping"),
			errors.WithMessage("failed to connect to Redis"),
			errors.WithErrorCode("redis/connection-failed"),
			errors.WithDetail(err.Error()),
		)
	}

	return &RedisClient{client: client}, nil
}
