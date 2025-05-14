package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	appPort := getEnv("APP_PORT", "8080")

	dbUser := getEnv("DB_USER", "")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "")
	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "")

	dbSource := ""
	if dbUser != "" && dbName != "" && dbHost != "" && dbPort != "" {
		dbSource = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, "disable",
		)
	}

	return Config{
		AppPort: appPort,

		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBUser:     dbUser,
		DBPassword: dbPass,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBSource:   dbSource, // constructed above

		RedisHost: getEnv("REDIS_HOST", ""),
		RedisPort: getEnv("REDIS_PORT", ""),

		ElasticURL: getEnv("ELASTIC_URL", "http://localhost:9200"),

		JwtSecret: getEnv("JWT_SECRET", "default-secret-key"),
		JwtExpiry: parseDuration(getEnv("JWT_EXPIRY", "15m")),
	}
}

// Helper: Get env var or fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Helper: Parse duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("Invalid duration: " + s)
	}
	return d
}
