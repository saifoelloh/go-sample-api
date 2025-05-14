package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/saifoelloh/ranger/internal/config"
	handler "github.com/saifoelloh/ranger/internal/handler"
	"github.com/saifoelloh/ranger/internal/middleware"
	repository "github.com/saifoelloh/ranger/internal/repositories"
	service "github.com/saifoelloh/ranger/internal/services"
	"github.com/saifoelloh/ranger/pkg/errors"
)

func initDB(cfg config.Config) *sqlx.DB {
	db, err := sqlx.Connect(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		errors.LogAndPanic(errors.InternalServerError(
			errors.WithScope("main"),
			errors.WithLocation("sqlx.Connect"),
			errors.WithMessage("failed to connect to database"),
			errors.WithErrorCode("db/connection-failed"),
			errors.WithDetail(err.Error()),
		))
	}

	log.Println("✅ PostgreSQL connected")
	return db
}

func initRedis(cfg config.Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
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

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db := initDB(cfg)
	redisClient := initRedis(cfg)

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize Services
	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JwtSecret)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup Router
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	// Routes
	router.POST("/login", authHandler.Login)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Run Server
	router.Run(":" + cfg.AppPort)
}
