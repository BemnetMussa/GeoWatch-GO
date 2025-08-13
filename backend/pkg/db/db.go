// pkg/db/db.go

package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB establishes a connection pool to the PostgreSQL database.
// It reads the connection string from an environment variable.
func ConnectDB() (*pgxpool.Pool, error) {
	// Example DSN: "postgres://username:password@localhost:5432/database_name"
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// pgxpool.New is the recommended way to create a connection pool.
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to verify the connection.
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close() // Close the pool if ping fails
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	fmt.Println("INFO: Successfully connected to PostgreSQL database.")
	return pool, nil
}