// shortener/db.go
package shortener

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a global database connection pool
var DB *pgxpool.Pool

// ConnectDB initializes the database connection using a URL from an environment variable.
func ConnectDB() error {
	databaseURL := os.Getenv("DATABASE_URL") // Neon database URL
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Create a new connection pool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create database connection pool: %v", err)
	}

	// Assign the pool to the global variable
	DB = pool
	log.Println("Connected to the database successfully!")
	return nil
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
