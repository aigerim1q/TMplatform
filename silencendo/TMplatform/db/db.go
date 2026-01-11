package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

// Config holds database configuration loaded from environment variables.
type Config struct {
	Driver string
	DSN    string
}

// LoadConfig constructs Config with safe defaults.
func LoadConfig() Config {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Durable file-based SQLite with foreign keys and busy timeout.
		dsn = "file:chatbot.db?_busy_timeout=5000&_foreign_keys=1"
	}

	return Config{Driver: driver, DSN: dsn}
}

// Connect opens a database connection and pings it.
func Connect(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

// WithTx runs fn within a transaction, committing on success and rolling back on failure.
func WithTx(ctx context.Context, dbConn *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error (%v) after %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
