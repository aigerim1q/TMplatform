package repository

import (
	"context"
	"database/sql"
	"fmt"

	"silencendo/models"
)

// TaskRepository manages tasks.
type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, tx *sql.Tx, task models.Task) (models.Task, error) {
	if tx == nil {
		return models.Task{}, fmt.Errorf("transaction required")
	}

	_, err := tx.ExecContext(ctx, `INSERT INTO tasks (id, project_id, stage_id, title, normalized_title, description, status, priority, assignee_id, due_date)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.ProjectID, task.StageID, task.Title, task.NormalizedTitle, task.Description, task.Status, task.Priority, task.AssigneeID, task.DueDate)
	if err != nil {
		return models.Task{}, fmt.Errorf("create task: %w", err)
	}
	return r.GetByID(ctx, tx, task.ID)
}

func (r *TaskRepository) GetByID(ctx context.Context, tx *sql.Tx, id string) (models.Task, error) {
	if tx == nil {
		return models.Task{}, fmt.Errorf("transaction required")
	}
	var t models.Task
	var stageID sql.NullString
	var assigneeID sql.NullString
	var due sql.NullTime

	err := tx.QueryRowContext(ctx, `SELECT id, project_id, stage_id, title, normalized_title, description, status, priority, assignee_id, due_date, created_at, updated_at FROM tasks WHERE id=?`, id).Scan(
		&t.ID, &t.ProjectID, &stageID, &t.Title, &t.NormalizedTitle, &t.Description, &t.Status, &t.Priority, &assigneeID, &due, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrNotFound
		}
		return models.Task{}, fmt.Errorf("get task: %w", err)
	}
	if stageID.Valid {
		val := stageID.String
		t.StageID = &val
	}
	if assigneeID.Valid {
		val := assigneeID.String
		t.AssigneeID = &val
	}
	if due.Valid {
		t.DueDate = &due.Time
	}
	return t, nil
}

func (r *TaskRepository) GetByProjectAndNormalizedTitle(ctx context.Context, tx *sql.Tx, projectID, normalized string) (models.Task, error) {
	if tx == nil {
		return models.Task{}, fmt.Errorf("transaction required")
	}
	var t models.Task
	var stageID sql.NullString
	var assigneeID sql.NullString
	var due sql.NullTime
	err := tx.QueryRowContext(ctx, `SELECT id, project_id, stage_id, title, normalized_title, description, status, priority, assignee_id, due_date, created_at, updated_at
        FROM tasks WHERE project_id=? AND normalized_title=?`, projectID, normalized).Scan(
		&t.ID, &t.ProjectID, &stageID, &t.Title, &t.NormalizedTitle, &t.Description, &t.Status, &t.Priority, &assigneeID, &due, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrNotFound
		}
		return models.Task{}, fmt.Errorf("get task by title: %w", err)
	}
	if stageID.Valid {
		val := stageID.String
		t.StageID = &val
	}
	if assigneeID.Valid {
		val := assigneeID.String
		t.AssigneeID = &val
	}
	if due.Valid {
		t.DueDate = &due.Time
	}
	return t, nil
}

func (r *TaskRepository) UpdateAssignee(ctx context.Context, tx *sql.Tx, taskID string, assigneeID *string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}
	_, err := tx.ExecContext(ctx, `UPDATE tasks SET assignee_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, assigneeID, taskID)
	if err != nil {
		return fmt.Errorf("update assignee: %w", err)
	}
	return nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, taskID, status string) error {
	if tx == nil {
		return fmt.Errorf("transaction required")
	}
	_, err := tx.ExecContext(ctx, `UPDATE tasks SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, taskID)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	return nil
}

func (r *TaskRepository) ListByProject(ctx context.Context, tx *sql.Tx, projectID string) ([]models.Task, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction required")
	}

	rows, err := tx.QueryContext(ctx, `SELECT id, project_id, stage_id, title, normalized_title, description, status, priority, assignee_id, due_date, created_at, updated_at
        FROM tasks WHERE project_id=? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var stageID sql.NullString
		var assigneeID sql.NullString
		var due sql.NullTime
		if err := rows.Scan(&t.ID, &t.ProjectID, &stageID, &t.Title, &t.NormalizedTitle, &t.Description, &t.Status, &t.Priority, &assigneeID, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		if stageID.Valid {
			val := stageID.String
			t.StageID = &val
		}
		if assigneeID.Valid {
			val := assigneeID.String
			t.AssigneeID = &val
		}
		if due.Valid {
			t.DueDate = &due.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) FindForUserByTitle(ctx context.Context, tx *sql.Tx, userID, normalized string) (models.Task, error) {
	if tx == nil {
		return models.Task{}, fmt.Errorf("transaction required")
	}

	var t models.Task
	var stageID sql.NullString
	var assigneeID sql.NullString
	var due sql.NullTime

	err := tx.QueryRowContext(ctx, `SELECT t.id, t.project_id, t.stage_id, t.title, t.normalized_title, t.description, t.status, t.priority, t.assignee_id, t.due_date, t.created_at, t.updated_at
            FROM tasks t
            JOIN projects p ON p.id = t.project_id
            JOIN project_members pm ON pm.project_id = p.id
            WHERE pm.user_id = ? AND t.normalized_title = ?
            LIMIT 1`, userID, normalized).Scan(
		&t.ID, &t.ProjectID, &stageID, &t.Title, &t.NormalizedTitle, &t.Description, &t.Status, &t.Priority, &assigneeID, &due, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrNotFound
		}
		return models.Task{}, fmt.Errorf("find task by title: %w", err)
	}
	if stageID.Valid {
		val := stageID.String
		t.StageID = &val
	}
	if assigneeID.Valid {
		val := assigneeID.String
		t.AssigneeID = &val
	}
	if due.Valid {
		t.DueDate = &due.Time
	}
	return t, nil
}
