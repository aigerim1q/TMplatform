package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"silencendo/models"
)

// ProjectRepository handles project persistence.
type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, tx *sql.Tx, p models.Project) (models.Project, error) {
	if tx == nil {
		return models.Project{}, fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `INSERT INTO projects (id, owner_id, title, description, status, normalized_title)
        VALUES (?, ?, ?, ?, ?, ?)`, p.ID, p.OwnerID, p.Title, p.Description, p.Status, p.NormalizedTitle)
	if err != nil {
		return models.Project{}, fmt.Errorf("create project: %w", err)
	}
	return r.GetByID(ctx, tx, p.ID)
}

func (r *ProjectRepository) GetByID(ctx context.Context, tx *sql.Tx, id string) (models.Project, error) {
	if tx == nil {
		return models.Project{}, fmt.Errorf("transaction required")
	}

	var p models.Project
	err := tx.QueryRowContext(ctx, `SELECT id, owner_id, title, description, status, normalized_title, deleted_at, created_at, updated_at
        FROM projects WHERE id=? AND deleted_at IS NULL`, id).Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.NormalizedTitle, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, fmt.Errorf("get project: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) GetByOwnerAndNormalizedTitle(ctx context.Context, tx *sql.Tx, ownerID, normalized string) (models.Project, error) {
	if tx == nil {
		return models.Project{}, fmt.Errorf("transaction required")
	}

	var p models.Project
	err := tx.QueryRowContext(ctx, `SELECT id, owner_id, title, description, status, normalized_title, deleted_at, created_at, updated_at
        FROM projects WHERE owner_id=? AND normalized_title=? AND deleted_at IS NULL`, ownerID, normalized).Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.NormalizedTitle, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, fmt.Errorf("get project by title: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) ListByUser(ctx context.Context, tx *sql.Tx, userID, search string, limit, offset int) ([]models.Project, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction required")
	}

	// Basic search by title using LIKE on normalized title.
	like := "%"
	if search != "" {
		like = "%" + strings.ToLower(strings.TrimSpace(search)) + "%"
	}

	rows, err := tx.QueryContext(ctx, `SELECT DISTINCT p.id, p.owner_id, p.title, p.description, p.status, p.normalized_title, p.deleted_at, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = ? AND p.deleted_at IS NULL AND (? = '%' OR p.normalized_title LIKE ?)
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, userID, like, like, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.NormalizedTitle, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// ListActiveByUser returns projects for the user that are currently marked active.
func (r *ProjectRepository) ListActiveByUser(ctx context.Context, tx *sql.Tx, userID string) ([]models.Project, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction required")
	}

	rows, err := tx.QueryContext(ctx, `SELECT DISTINCT p.id, p.owner_id, p.title, p.description, p.status, p.normalized_title, p.deleted_at, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = ? AND p.status = 'active' AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list active projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.NormalizedTitle, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan active project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// UpdateStatus sets the status for a project.
func (r *ProjectRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, projectID, status string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `UPDATE projects SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, projectID)
	if err != nil {
		return fmt.Errorf("update project status: %w", err)
	}
	return nil
}

// SoftDelete marks a project as deleted without removing rows.
func (r *ProjectRepository) SoftDelete(ctx context.Context, tx *sql.Tx, projectID string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `UPDATE projects SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, projectID)
	if err != nil {
		return fmt.Errorf("soft delete project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) FindByUserAndNormalizedTitle(ctx context.Context, tx *sql.Tx, userID, normalized string) (models.Project, error) {
	if tx == nil {
		return models.Project{}, fmt.Errorf("transaction required")
	}

	var p models.Project
	err := tx.QueryRowContext(ctx, `SELECT p.id, p.owner_id, p.title, p.description, p.status, p.normalized_title, p.deleted_at, p.created_at, p.updated_at
        FROM projects p
        JOIN project_members pm ON pm.project_id = p.id
        WHERE pm.user_id = ? AND p.normalized_title = ? AND p.deleted_at IS NULL
        LIMIT 1`, userID, normalized).Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.NormalizedTitle, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, fmt.Errorf("find project by user/title: %w", err)
	}
	return p, nil
}
