package repository

import (
	"context"
	"database/sql"
	"fmt"

	"silencendo/models"
)

// StageRepository handles persistence for stages.
type StageRepository struct {
	db *sql.DB
}

func NewStageRepository(db *sql.DB) *StageRepository {
	return &StageRepository{db: db}
}

func (r *StageRepository) Create(ctx context.Context, tx *sql.Tx, stage models.Stage) (models.Stage, error) {
	if tx == nil {
		return models.Stage{}, fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `INSERT INTO stages (id, project_id, title, normalized_title, order_index, responsible_id)
        VALUES (?, ?, ?, ?, ?, ?)`, stage.ID, stage.ProjectID, stage.Title, stage.NormalizedTitle, stage.OrderIndex, stage.ResponsibleID)
	if err != nil {
		return models.Stage{}, fmt.Errorf("create stage: %w", err)
	}
	return r.GetByID(ctx, tx, stage.ID)
}

func (r *StageRepository) GetByID(ctx context.Context, tx *sql.Tx, id string) (models.Stage, error) {
	if tx == nil {
		return models.Stage{}, fmt.Errorf("transaction required")
	}

	var s models.Stage
	var responsibleID sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT id, project_id, title, normalized_title, order_index, responsible_id, created_at, updated_at FROM stages WHERE id=?`, id).Scan(
		&s.ID, &s.ProjectID, &s.Title, &s.NormalizedTitle, &s.OrderIndex, &responsibleID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Stage{}, ErrNotFound
		}
		return models.Stage{}, fmt.Errorf("get stage: %w", err)
	}
	if responsibleID.Valid {
		val := responsibleID.String
		s.ResponsibleID = &val
	}
	return s, nil
}

func (r *StageRepository) GetByProjectAndNormalizedTitle(ctx context.Context, tx *sql.Tx, projectID, normalized string) (models.Stage, error) {
	if tx == nil {
		return models.Stage{}, fmt.Errorf("transaction required")
	}
	var s models.Stage
	var responsibleID sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT id, project_id, title, normalized_title, order_index, responsible_id, created_at, updated_at FROM stages
        WHERE project_id=? AND normalized_title=?`, projectID, normalized).Scan(
		&s.ID, &s.ProjectID, &s.Title, &s.NormalizedTitle, &s.OrderIndex, &responsibleID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Stage{}, ErrNotFound
		}
		return models.Stage{}, fmt.Errorf("get stage by title: %w", err)
	}
	if responsibleID.Valid {
		val := responsibleID.String
		s.ResponsibleID = &val
	}
	return s, nil
}

func (r *StageRepository) ListByProject(ctx context.Context, tx *sql.Tx, projectID string) ([]models.Stage, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction required")
	}

	rows, err := tx.QueryContext(ctx, `SELECT id, project_id, title, normalized_title, order_index, responsible_id, created_at, updated_at FROM stages
        WHERE project_id=? ORDER BY order_index ASC, created_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list stages: %w", err)
	}
	defer rows.Close()

	var stages []models.Stage
	for rows.Next() {
		var s models.Stage
		var responsibleID sql.NullString
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Title, &s.NormalizedTitle, &s.OrderIndex, &responsibleID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan stage: %w", err)
		}
		if responsibleID.Valid {
			val := responsibleID.String
			s.ResponsibleID = &val
		}
		stages = append(stages, s)
	}
	return stages, nil
}

func (r *StageRepository) UpdateResponsible(ctx context.Context, tx *sql.Tx, stageID string, responsibleID *string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}
	_, err := tx.ExecContext(ctx, `UPDATE stages SET responsible_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, responsibleID, stageID)
	if err != nil {
		return fmt.Errorf("update stage responsible: %w", err)
	}
	return nil
}

func (r *StageRepository) FindForUserByTitle(ctx context.Context, tx *sql.Tx, userID, normalized string) (models.Stage, error) {
	if tx == nil {
		return models.Stage{}, fmt.Errorf("transaction required")
	}

	var s models.Stage
	var responsibleID sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT s.id, s.project_id, s.title, s.normalized_title, s.order_index, s.responsible_id, s.created_at, s.updated_at
        FROM stages s
        JOIN projects p ON p.id = s.project_id
        JOIN project_members pm ON pm.project_id = p.id
        WHERE pm.user_id = ? AND s.normalized_title = ?
        LIMIT 1`, userID, normalized).Scan(
		&s.ID, &s.ProjectID, &s.Title, &s.NormalizedTitle, &s.OrderIndex, &responsibleID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Stage{}, ErrNotFound
		}
		return models.Stage{}, fmt.Errorf("find stage by title: %w", err)
	}
	if responsibleID.Valid {
		val := responsibleID.String
		s.ResponsibleID = &val
	}
	return s, nil
}
