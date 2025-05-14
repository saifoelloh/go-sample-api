package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saifoelloh/ranger/pkg/errors"
)

const (
	AccessTokenKey   = "token:access:%s" // %s = userID
	UserIDByTokenKey = "token:user:%s"   // %s = accessToken
)

type TokenRepository struct {
	client *RedisClient
}

func NewTokenRepository(client *RedisClient) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) SetAccessToken(ctx context.Context, userID, accessToken string, ttl time.Duration) error {
	accessTokenKey := fmt.Sprintf(AccessTokenKey, userID)
	userIDKey := fmt.Sprintf(UserIDByTokenKey, accessToken)

	pipe := r.client.Client.Pipeline()
	pipe.Set(ctx, accessTokenKey, userID, ttl)
	pipe.Set(ctx, userIDKey, userID, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errors.InternalServerError(
			errors.WithScope("TokenRepository"),
			errors.WithLocation("SetAccessToken.Exec"),
			errors.WithMessage("failed to store access token in Redis"),
			errors.WithErrorCode("redis/set-token-failed"),
		)
	}

	return nil
}

func (r *TokenRepository) GetUserIDFromToken(ctx context.Context, accessToken string) (string, error) {
	key := fmt.Sprintf(UserIDByTokenKey, accessToken)
	userID, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.Unauthorized(
				errors.WithScope("TokenRepository"),
				errors.WithLocation("GetUserIDFromToken.NotFound"),
				errors.WithMessage("token not found or expired"),
				errors.WithErrorCode("auth/token-invalid-or-expired"),
			)
		}
		return "", errors.InternalServerError(
			errors.WithScope("TokenRepository"),
			errors.WithLocation("GetUserIDFromToken.RedisError"),
			errors.WithMessage("failed to fetch user from token"),
			errors.WithErrorCode("auth/token-validation-failed"),
		)
	}
	return userID, nil
}
