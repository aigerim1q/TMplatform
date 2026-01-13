package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// IdempotencyRepository stores processed keys.
type IdempotencyRepository struct {
	db *sql.DB
}

func NewIdempotencyRepository(db *sql.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Record(ctx context.Context, tx *sql.Tx, key, userID, intent, resourceType, resourceID string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO idempotency_keys (key, user_id, intent, resource_type, resource_id)
        VALUES (?, ?, ?, ?, ?)`, key, userID, intent, resourceType, resourceID)
	if err != nil {
		return fmt.Errorf("record idempotency: %w", err)
	}
	return nil
}

// Lookup returns resource info if the key already exists.
func (r *IdempotencyRepository) Lookup(ctx context.Context, tx *sql.Tx, key, userID string) (resourceType string, resourceID string, found bool, err error) {
	if tx == nil {
		return "", "", false, fmt.Errorf("transaction required")
	}
	err = tx.QueryRowContext(ctx, `SELECT resource_type, resource_id FROM idempotency_keys WHERE key=? AND user_id=?`, key, userID).Scan(&resourceType, &resourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("lookup idempotency: %w", err)
	}
	return resourceType, resourceID, true, nil
}
