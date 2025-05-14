package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/saifoelloh/ranger/internal/config"
	handler "github.com/saifoelloh/ranger/internal/handler"
	"github.com/saifoelloh/ranger/internal/middleware"
	"github.com/saifoelloh/ranger/internal/redis"
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

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db := initDB(cfg)
	rdb := redis.InitRedis(cfg)

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Redis
	limiterCfg := redis.RateLimiterConfig{
		MaxAttempts:     3,
		DelayPerAttempt: 10 * time.Second,
		LockoutDuration: 10 * time.Minute,
	}
	redisClient := redis.NewRedisClient(rdb)
	rateLimiterRepo := redis.NewRateLimiterRepository(redisClient, limiterCfg)
	tokenCacheRepo := redis.NewTokenRepository(redisClient)

	// Initialize Services
	authService := service.NewAuthService(cfg, userRepo, sessionRepo, rateLimiterRepo, tokenCacheRepo)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup Router
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		if c.Request.URL.Port() != cfg.AppPort {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	// Routes
	router.POST("/login", authHandler.Login)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Run Server
	router.Run(":" + cfg.AppPort)
}
