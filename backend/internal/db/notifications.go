package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
	EntityType *string   `json:"entity_type,omitempty"`
	EntityID   *int      `json:"entity_id,omitempty"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

func ListNotifications(ctx context.Context, pool *pgxpool.Pool, userID int, limit int, offset int) ([]Notification, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	const q = `
SELECT id, user_id, type, message, entity_type, entity_id, is_read, created_at
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
`
	rows, err := pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Notification, 0)
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.EntityType, &n.EntityID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func MarkNotificationRead(ctx context.Context, pool *pgxpool.Pool, userID int, notificationID int) (bool, error) {
	const q = `
UPDATE notifications
SET is_read = TRUE
WHERE id = $1 AND user_id = $2;
`
	ct, err := pool.Exec(ctx, q, notificationID, userID)
	if err != nil {
		return false, err
	}
	return ct.RowsAffected() > 0, nil
}
