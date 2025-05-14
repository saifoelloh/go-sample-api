package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/saifoelloh/ranger/internal/config"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{Client: client}
}

func InitRedis(cfg config.Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	fmt.Println(redisAddr)
	if redisAddr == "" {
		log.Println("🟡 Redis not configured")
		return nil
	}

	client := redis.NewClient(&redis.Options{Addr: redisAddr})

	if err := client.Ping(context.Background()).Err(); err != nil {
		errors.LogAndPanic(errors.InternalServerError(
			errors.WithScope("main"),
			errors.WithLocation("redis.Ping"),
			errors.WithMessage("failed to connect to Redis"),
			errors.WithErrorCode("redis/connection-failed"),
			errors.WithDetail(err.Error()),
		))
	}

	log.Println("✅ Redis connected")
	return client
}
