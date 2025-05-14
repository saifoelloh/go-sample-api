package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	sslMode      = "disable"
	maxAttempts  = 3
	retryDelay   = 5 * time.Second
	maxOpenConns = 10
	maxIdleConns = 10
	connMaxLife  = 2 * time.Hour
)

// Connect initializes the DB connection using sqlx with retry and pooling configuration
func Connect(connStr string) (*sqlx.DB, error) {
	var err error

	fmt.Println(connStr)
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var db *sqlx.DB
		db, err = sqlx.Connect("postgres", connStr)
		if err == nil {
			log.Println("[DB] Connection established.")
			configureConnectionPool(db)
			return db, nil
		}

		log.Printf("[DB] Attempt %d: failed to connect: %v", attempt, err)
		time.Sleep(retryDelay)
	}

	return nil, errors.New("unable to establish DB connection after multiple attempts")
}

// configureConnectionPool sets connection pooling parameters
func configureConnectionPool(db *sqlx.DB) {
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLife)
}
