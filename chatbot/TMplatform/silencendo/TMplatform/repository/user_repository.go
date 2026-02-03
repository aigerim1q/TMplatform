package repository

import (
	"context"
	"database/sql"
	"fmt"

	"silencendo/models"
)

// UserRepository manages persistence for users.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Upsert(ctx context.Context, tx *sql.Tx, user models.User) (models.User, error) {
	if tx == nil {
		return models.User{}, fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `INSERT INTO users (id, name, email)
        VALUES (?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            name=excluded.name,
            email=excluded.email,
            updated_at=CURRENT_TIMESTAMP`, user.ID, user.Name, user.Email)
	if err != nil {
		return models.User{}, fmt.Errorf("upsert user: %w", err)
	}

	return r.GetByID(ctx, tx, user.ID)
}

func (r *UserRepository) GetByID(ctx context.Context, tx *sql.Tx, id string) (models.User, error) {
	if tx == nil {
		return models.User{}, fmt.Errorf("transaction required")
	}

	var u models.User
	var email sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT id, name, email, created_at, updated_at FROM users WHERE id=?`, id).Scan(
		&u.ID, &u.Name, &email, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}
	if email.Valid {
		u.Email = &email.String
	}
	return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (models.User, error) {
	if tx == nil {
		return models.User{}, fmt.Errorf("transaction required")
	}

	var u models.User
	var emailPtr sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT id, name, email, created_at, updated_at FROM users WHERE email=?`, email).Scan(
		&u.ID, &u.Name, &emailPtr, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}
	if emailPtr.Valid {
		u.Email = &emailPtr.String
	}
	return u, nil
}
