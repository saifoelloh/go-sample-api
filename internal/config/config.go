package config

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	AppPort string

	DBDriver   string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBSource   string // constructed from other fields

	RedisHost string
	RedisPort string

	ElasticURL string

	JwtSecret string
	JwtExpiry time.Duration
	JwtIssuer string
}

var (
	DB  *sqlx.DB
	RDB *redis.Client
)
