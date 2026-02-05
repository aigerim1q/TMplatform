package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int        `json:"id"`
	StageID     int        `json:"stage_id"`
	ParentID    *int       `json:"parent_id,omitempty"`
	AssigneeID  *int       `json:"assignee_id,omitempty"`
	AuthorID    int        `json:"author_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"` // todo, in_progress, done, review
	Priority    int        `json:"priority"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Joined fields
	AssigneeName  *string `json:"assignee_name,omitempty"`
	AuthorName    string  `json:"author_name,omitempty"`
	FilesCount    int     `json:"files_count"`
	SubtasksCount int     `json:"subtasks_count"`
}

type TaskSummary struct {
	ID              int        `json:"id"`
	StageID         int        `json:"stage_id"`
	StageTitle      string     `json:"stage_title"`
	ProjectID       int        `json:"project_id"`
	ProjectName     string     `json:"project_name"`
	ProjectImageURL *string    `json:"project_image_url,omitempty"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	AssigneeID      *int       `json:"assignee_id,omitempty"`
	AssigneeName    *string    `json:"assignee_name,omitempty"`
	AuthorName      string     `json:"author_name"`
}

func CreateTask(ctx context.Context, db *pgxpool.Pool, t *Task) error {
	q := `
INSERT INTO tasks (stage_id, parent_id, assignee_id, author_id, title, description, status, priority, start_date, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at, updated_at
`
	return db.QueryRow(ctx, q,
		t.StageID, t.ParentID, t.AssigneeID, t.AuthorID, t.Title, t.Description, t.Status, t.Priority, t.StartDate, t.DueDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func GetTaskByID(ctx context.Context, db *pgxpool.Pool, taskID int) (*Task, error) {
	q := `
SELECT t.id, t.stage_id, t.parent_id, t.assignee_id, t.author_id, t.title, t.description, t.status, t.priority, t.start_date, t.due_date, t.created_at, t.updated_at,
	   COALESCE(u1.display_name, u1.email) as assignee_name, COALESCE(u2.display_name, u2.email) as author_name,
       (SELECT COUNT(*) FROM files f WHERE f.task_id = t.id) as files_count,
       (SELECT COUNT(*) FROM tasks st WHERE st.parent_id = t.id) as subtasks_count
FROM tasks t
LEFT JOIN users u1 ON t.assignee_id = u1.id
JOIN users u2 ON t.author_id = u2.id
WHERE t.id = $1
`
	var t Task
	err := db.QueryRow(ctx, q, taskID).Scan(
		&t.ID, &t.StageID, &t.ParentID, &t.AssigneeID, &t.AuthorID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.StartDate, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		&t.AssigneeName, &t.AuthorName, &t.FilesCount, &t.SubtasksCount,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ListTasksByStage(ctx context.Context, db *pgxpool.Pool, stageID int) ([]Task, error) {
	q := `
SELECT t.id, t.stage_id, t.parent_id, t.assignee_id, t.author_id, t.title, t.description, t.status, t.priority, t.start_date, t.due_date, t.created_at, t.updated_at,
	   COALESCE(u1.display_name, u1.email) as assignee_name, COALESCE(u2.display_name, u2.email) as author_name,
       (SELECT COUNT(*) FROM files f WHERE f.task_id = t.id) as files_count,
       (SELECT COUNT(*) FROM tasks st WHERE st.parent_id = t.id) as subtasks_count
FROM tasks t
LEFT JOIN users u1 ON t.assignee_id = u1.id
JOIN users u2 ON t.author_id = u2.id
WHERE t.stage_id = $1 AND t.parent_id IS NULL
ORDER BY t.priority DESC, t.created_at ASC
`
	return scanTasks(ctx, db, q, stageID)
}

func ListSubtasks(ctx context.Context, db *pgxpool.Pool, parentID int) ([]Task, error) {
	q := `
SELECT t.id, t.stage_id, t.parent_id, t.assignee_id, t.author_id, t.title, t.description, t.status, t.priority, t.start_date, t.due_date, t.created_at, t.updated_at,
	   COALESCE(u1.display_name, u1.email) as assignee_name, COALESCE(u2.display_name, u2.email) as author_name,
       (SELECT COUNT(*) FROM files f WHERE f.task_id = t.id) as files_count,
       (SELECT COUNT(*) FROM tasks st WHERE st.parent_id = t.id) as subtasks_count
FROM tasks t
LEFT JOIN users u1 ON t.assignee_id = u1.id
JOIN users u2 ON t.author_id = u2.id
WHERE t.parent_id = $1
ORDER BY t.priority DESC, t.created_at ASC
`
	return scanTasks(ctx, db, q, parentID)
}

func UpdateTask(ctx context.Context, db *pgxpool.Pool, t *Task) error {
	q := `
UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, start_date = $5, due_date = $6, assignee_id = $7, updated_at = now()
WHERE id = $8
`
	_, err := db.Exec(ctx, q, t.Title, t.Description, t.Status, t.Priority, t.StartDate, t.DueDate, t.AssigneeID, t.ID)
	return err
}

func scanTasks(ctx context.Context, db *pgxpool.Pool, query string, args ...any) ([]Task, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID, &t.StageID, &t.ParentID, &t.AssigneeID, &t.AuthorID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.StartDate, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
			&t.AssigneeName, &t.AuthorName, &t.FilesCount, &t.SubtasksCount,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func ListTasksForDashboard(ctx context.Context, db *pgxpool.Pool, orgID int, userID int, scope string, limit int, projectID int, mine bool) ([]TaskSummary, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	where := "p.org_id = $1 AND p.status <> 'deleted' AND t.parent_id IS NULL"
	args := []any{orgID}
	argPos := 2

	if projectID > 0 {
		where = fmt.Sprintf("%s AND p.id = $%d", where, argPos)
		args = append(args, projectID)
		argPos++
	}

	switch scope {
	case "urgent":
		// Срочные: задачи с дедлайном <=5 дней или задачи проектов с дедлайном <=5 дней
		where = fmt.Sprintf("%s AND (t.assignee_id = $%d OR t.author_id = $%d) AND (t.due_date IS NOT NULL AND t.due_date < now() + interval '5 days' OR p.end_date IS NOT NULL AND p.end_date < now() + interval '5 days')", where, argPos, argPos)
		args = append(args, userID)
		argPos++
	case "subordinate":
		// Подчинённые: пользователи с manager_id = current user OR ответственные в проектах, где текущий пользователь владелец
		where = fmt.Sprintf(`%s AND t.assignee_id IN (
			SELECT id FROM users WHERE manager_id = $%d
			UNION
			SELECT pa.user_id FROM project_assignees pa
			JOIN projects p2 ON p2.id = pa.project_id
			WHERE p2.owner_id = $%d
		)`, where, argPos, argPos)
		args = append(args, userID)
		argPos++
	default:
		if mine {
			where = fmt.Sprintf("%s AND (t.assignee_id = $%d OR t.author_id = $%d)", where, argPos, argPos+1)
			args = append(args, userID, userID)
			argPos += 2
		} else {
			where = fmt.Sprintf("%s AND t.assignee_id = $%d", where, argPos)
			args = append(args, userID)
			argPos++
		}
	}

	q := fmt.Sprintf(`
SELECT t.id, t.stage_id, s.title, p.id, p.name, p.image_url,
       t.title, t.description, t.status, t.priority, t.due_date,
	   t.assignee_id, COALESCE(u1.display_name, u1.email) as assignee_name, COALESCE(u2.display_name, u2.email) as author_name
FROM tasks t
JOIN stages s ON t.stage_id = s.id
JOIN projects p ON s.project_id = p.id
LEFT JOIN users u1 ON t.assignee_id = u1.id
JOIN users u2 ON t.author_id = u2.id
WHERE %s
ORDER BY COALESCE(t.due_date, p.end_date) NULLS LAST, t.created_at DESC
LIMIT $%d
`, where, argPos)

	args = append(args, limit)
	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []TaskSummary{}
	for rows.Next() {
		var t TaskSummary
		if err := rows.Scan(
			&t.ID,
			&t.StageID,
			&t.StageTitle,
			&t.ProjectID,
			&t.ProjectName,
			&t.ProjectImageURL,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.AssigneeID,
			&t.AssigneeName,
			&t.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, nil
}
