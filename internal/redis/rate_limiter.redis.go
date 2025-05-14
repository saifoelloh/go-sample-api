package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saifoelloh/ranger/pkg/errors"
)

const (
	loginRateLimitKey = "rate-limit:login:%s"  // %s = ip
	apiRateLimitKey   = "rate-limit:api:%s:%s" // %s = userID, endpoint
)

type RateLimiterRepository struct {
	client *RedisClient
	cfg    RateLimiterConfig
}

type RateLimiterConfig struct {
	MaxAttempts     int
	DelayPerAttempt time.Duration
	LockoutDuration time.Duration
}

func NewRateLimiterRepository(client *RedisClient, cfg RateLimiterConfig) *RateLimiterRepository {
	return &RateLimiterRepository{
		client: client,
		cfg:    cfg,
	}
}

func (r *RateLimiterRepository) IsAllowed(ctx context.Context, ip string) error {
	key := fmt.Sprintf(loginRateLimitKey, ip)
	cmd := r.client.Client

	err := cmd.Watch(ctx, func(tx *redis.Tx) error {
		// Cek apakah key ada
		_, err := tx.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return errors.InternalServerError(
				errors.WithScope("RateLimiter"),
				errors.WithLocation("IsAllowed.Get"),
				errors.WithMessage("failed to get rate limit data"),
				errors.WithErrorCode("redis/get-error"),
			)
		}

		var attempts int64 = 0
		if err == nil {
			attempts, err = tx.Get(ctx, key).Int64()
			if err != nil {
				return errors.InternalServerError(
					errors.WithScope("RateLimiter"),
					errors.WithLocation("IsAllowed.Parse"),
					errors.WithMessage("failed to parse attempts count"),
					errors.WithErrorCode("redis/parse-error"),
				)
			}
		}

		if attempts >= int64(r.cfg.MaxAttempts) {
			// Lockout: extend expiration
			_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Expire(ctx, key, r.cfg.LockoutDuration)
				return nil
			})
			if err != nil {
				return errors.InternalServerError(
					errors.WithScope("RateLimiter"),
					errors.WithLocation("IsAllowed.Expire"),
					errors.WithMessage("failed to extend lockout"),
					errors.WithErrorCode("redis/expire-error"),
				)
			}
			return errors.TooManyRequests(
				errors.WithScope("RateLimiter"),
				errors.WithLocation("IsAllowed.AttemptsExceeded"),
				errors.WithMessage("too many login attempts. try again later"),
				errors.WithErrorCode("auth/too-many-attempts"),
			)
		}

		// Tambah attempts dan set TTL jika baru
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Incr(ctx, key)
			if attempts == 0 {
				pipe.Expire(ctx, key, r.cfg.DelayPerAttempt)
			}
			return nil
		})
		if err != nil {
			return errors.InternalServerError(
				errors.WithScope("RateLimiter"),
				errors.WithLocation("IsAllowed.Incr"),
				errors.WithMessage("failed to increment rate limit"),
				errors.WithErrorCode("redis/incr-error"),
			)
		}

		return nil
	}, key)

	if err != nil {
		// Return kalau memang error dari custom error handler di atas
		if e, ok := err.(*errors.Extension); ok {
			return e
		}
		return errors.InternalServerError(
			errors.WithScope("RateLimiter"),
			errors.WithLocation("IsAllowed.Watch"),
			errors.WithMessage("unexpected error during rate limiting"),
			errors.WithErrorCode("redis/watch-transaction-error"),
			errors.WithDetail(err.Error()),
		)
	}

	return nil
}

func (r *RateLimiterRepository) Reset(ctx context.Context, ip string) error {
	key := fmt.Sprintf(loginRateLimitKey, ip)
	err := r.client.Client.Del(ctx, key).Err()
	if err != nil {
		return errors.InternalServerError(
			errors.WithScope("RateLimiter"),
			errors.WithLocation("Reset.Del"),
			errors.WithMessage("failed to reset rate limit"),
			errors.WithErrorCode("redis/del-rate-limit-failed"),
		)
	}
	return nil
}
