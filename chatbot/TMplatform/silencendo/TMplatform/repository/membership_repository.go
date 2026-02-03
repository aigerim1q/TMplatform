package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// MembershipRepository manages project membership records.
type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(db *sql.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) AddIfMissing(ctx context.Context, tx *sql.Tx, projectID, userID, role string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO project_members (project_id, user_id, role)
        VALUES (?, ?, ?)
        ON CONFLICT(project_id, user_id) DO UPDATE SET role=excluded.role`, projectID, userID, role)
	if err != nil {
		return fmt.Errorf("add membership: %w", err)
	}
	return nil
}

func (r *MembershipRepository) IsMember(ctx context.Context, tx *sql.Tx, projectID, userID string) (bool, error) {
	if tx == nil {
		return false, fmt.Errorf("transaction required")
	}
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM project_members WHERE project_id=? AND user_id=?`, projectID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check membership: %w", err)
	}
	return count > 0, nil
}
