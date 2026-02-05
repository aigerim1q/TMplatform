package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID              int        `json:"id"`
	OrgID           int        `json:"org_id"`
	OwnerID         int        `json:"owner_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Status          string     `json:"status"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Priority        int        `json:"priority"`
	ImageURL        *string    `json:"image_url,omitempty"`
	BudgetAllocated float64    `json:"budget_allocated"`
	BudgetSpent     float64    `json:"budget_spent"`
	BudgetCurrency  string     `json:"budget_currency"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ProgressPercent int        `json:"progress_percent"`
	StagesCount     int        `json:"stages_count"`
	DoneStagesCount int        `json:"done_stages_count"`
	Assignees       []Assignee `json:"assignees"`
}

type Assignee struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type Stage struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	OrderIndex  int        `json:"order_index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProjectInput struct {
	OrgID           int
	OwnerID         int
	Name            string
	Description     *string
	Status          string
	StartDate       *time.Time
	EndDate         *time.Time
	Priority        int
	ImageURL        *string
	BudgetAllocated float64
	BudgetSpent     float64
	BudgetCurrency  string
	AssigneeIDs     []int
}

func ListProjects(ctx context.Context, pool *pgxpool.Pool, orgID int, userID int, mineOnly bool, limit int, offset int, sort string, order string) ([]Project, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	orderBy := projectOrderBy(sort, order)

	filter := ``
	args := []any{orgID, limit, offset}

	if mineOnly && userID > 0 {
		filter = ` AND (p.owner_id = $4 OR EXISTS (SELECT 1 FROM project_assignees pa WHERE pa.project_id = p.id AND pa.user_id = $4))`
		args = append(args, userID)
	}

	q := fmt.Sprintf(`
SELECT p.id, p.org_id, p.owner_id, p.name, p.description, p.status, p.start_date, p.end_date,
       p.priority, p.image_url, p.budget_allocated, p.budget_spent, p.budget_currency,
       p.created_at, p.updated_at,
       COALESCE(SUM(CASE WHEN s.status = 'done' THEN 1 ELSE 0 END), 0) AS done_count,
       COALESCE(COUNT(s.id), 0) AS total_count
FROM projects p
LEFT JOIN stages s ON s.project_id = p.id
WHERE p.org_id = $1 AND p.status <> 'deleted'%s
GROUP BY p.id
ORDER BY %s
LIMIT $2 OFFSET $3;`, filter, orderBy)

	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Project, 0)
	for rows.Next() {
		var p Project
		var doneCount int
		var totalCount int
		if err := rows.Scan(
			&p.ID,
			&p.OrgID,
			&p.OwnerID,
			&p.Name,
			&p.Description,
			&p.Status,
			&p.StartDate,
			&p.EndDate,
			&p.Priority,
			&p.ImageURL,
			&p.BudgetAllocated,
			&p.BudgetSpent,
			&p.BudgetCurrency,
			&p.CreatedAt,
			&p.UpdatedAt,
			&doneCount,
			&totalCount,
		); err != nil {
			return nil, err
		}
		p.StagesCount = totalCount
		p.DoneStagesCount = doneCount
		p.ProgressPercent = progressPercent(doneCount, totalCount)
		p.Assignees = []Assignee{}
		out = append(out, p)
	}
	return out, rows.Err()
}

func GetProject(ctx context.Context, pool *pgxpool.Pool, orgID int, projectID int) (*Project, error) {
	const q = `
SELECT p.id, p.org_id, p.owner_id, p.name, p.description, p.status, p.start_date, p.end_date,
       p.priority, p.image_url, p.budget_allocated, p.budget_spent, p.budget_currency,
       p.created_at, p.updated_at,
       COALESCE(SUM(CASE WHEN s.status = 'done' THEN 1 ELSE 0 END), 0) AS done_count,
       COALESCE(COUNT(s.id), 0) AS total_count
FROM projects p
LEFT JOIN stages s ON s.project_id = p.id
WHERE p.org_id = $1 AND p.id = $2 AND p.status <> 'deleted'
GROUP BY p.id
LIMIT 1;`

	var p Project
	var doneCount int
	var totalCount int
	err := pool.QueryRow(ctx, q, orgID, projectID).Scan(
		&p.ID,
		&p.OrgID,
		&p.OwnerID,
		&p.Name,
		&p.Description,
		&p.Status,
		&p.StartDate,
		&p.EndDate,
		&p.Priority,
		&p.ImageURL,
		&p.BudgetAllocated,
		&p.BudgetSpent,
		&p.BudgetCurrency,
		&p.CreatedAt,
		&p.UpdatedAt,
		&doneCount,
		&totalCount,
	)
	if err != nil {
		return nil, err
	}
	p.StagesCount = totalCount
	p.DoneStagesCount = doneCount
	p.ProgressPercent = progressPercent(doneCount, totalCount)

	assignees, err := listProjectAssignees(ctx, pool, projectID)
	if err != nil {
		return nil, err
	}
	p.Assignees = assignees

	return &p, nil
}

func CreateProject(ctx context.Context, pool *pgxpool.Pool, input CreateProjectInput) (*Project, error) {
	if input.Status == "" {
		input.Status = "draft"
	}
	if strings.TrimSpace(input.BudgetCurrency) == "" {
		input.BudgetCurrency = "KZT"
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var projectID int
	const insertProject = `
INSERT INTO projects (
  org_id, owner_id, name, description, status, start_date, end_date, priority,
  image_url, budget_allocated, budget_spent, budget_currency
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
RETURNING id;`

	err = tx.QueryRow(ctx, insertProject,
		input.OrgID,
		input.OwnerID,
		input.Name,
		input.Description,
		input.Status,
		input.StartDate,
		input.EndDate,
		input.Priority,
		input.ImageURL,
		input.BudgetAllocated,
		input.BudgetSpent,
		input.BudgetCurrency,
	).Scan(&projectID)
	if err != nil {
		return nil, err
	}

	if len(input.AssigneeIDs) > 0 {
		const insertAssignee = `
INSERT INTO project_assignees (project_id, user_id)
VALUES ($1,$2)
ON CONFLICT DO NOTHING;`
		for _, uid := range input.AssigneeIDs {
			if uid <= 0 {
				continue
			}
			if _, err := tx.Exec(ctx, insertAssignee, projectID, uid); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return GetProject(ctx, pool, input.OrgID, projectID)
}

func SoftDeleteProject(ctx context.Context, pool *pgxpool.Pool, orgID int, projectID int) (bool, error) {
	// Hard delete to avoid status check constraint issues; cascades will clean stages/tasks/files.
	res, err := pool.Exec(ctx,
		`DELETE FROM projects WHERE id = $1 AND org_id = $2`,
		projectID, orgID,
	)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func ListStages(ctx context.Context, pool *pgxpool.Pool, projectID int, sort string, order string) ([]Stage, error) {
	orderBy := stageOrderBy(sort, order)

	q := fmt.Sprintf(`
SELECT id, project_id, title, description, status, start_date, end_date, order_index, created_at, updated_at
FROM stages
WHERE project_id = $1
ORDER BY %s;`, orderBy)

	rows, err := pool.Query(ctx, q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Stage, 0)
	for rows.Next() {
		var s Stage
		if err := rows.Scan(
			&s.ID,
			&s.ProjectID,
			&s.Title,
			&s.Description,
			&s.Status,
			&s.StartDate,
			&s.EndDate,
			&s.OrderIndex,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func CreateStage(ctx context.Context, pool *pgxpool.Pool, projectID int, title string, description *string) (*Stage, error) {
	q := `
INSERT INTO stages (project_id, title, description, order_index)
VALUES ($1, $2, $3, COALESCE((SELECT MAX(order_index) + 1 FROM stages WHERE project_id = $1), 1))
RETURNING id, project_id, title, description, status, start_date, end_date, order_index, created_at, updated_at
`

	var s Stage
	if err := pool.QueryRow(ctx, q, projectID, title, description).Scan(
		&s.ID,
		&s.ProjectID,
		&s.Title,
		&s.Description,
		&s.Status,
		&s.StartDate,
		&s.EndDate,
		&s.OrderIndex,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &s, nil
}

func UpdateStageStatus(ctx context.Context, pool *pgxpool.Pool, stageID int, status string) (*Stage, error) {
	const q = `
UPDATE stages
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, project_id, title, description, status, start_date, end_date, order_index, created_at, updated_at;`

	var s Stage
	if err := pool.QueryRow(ctx, q, stageID, status).Scan(
		&s.ID,
		&s.ProjectID,
		&s.Title,
		&s.Description,
		&s.Status,
		&s.StartDate,
		&s.EndDate,
		&s.OrderIndex,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func listProjectAssignees(ctx context.Context, pool *pgxpool.Pool, projectID int) ([]Assignee, error) {
	const q = `
SELECT u.id, u.email, COALESCE(pa.role, u.role) AS role
FROM project_assignees pa
JOIN users u ON u.id = pa.user_id
WHERE pa.project_id = $1
ORDER BY u.email ASC;`

	rows, err := pool.Query(ctx, q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Assignee, 0)
	for rows.Next() {
		var a Assignee
		if err := rows.Scan(&a.ID, &a.Email, &a.Role); err != nil {
			return nil, err
		}
		a.Name = a.Email
		out = append(out, a)
	}
	return out, rows.Err()
}

func ListProjectAssignees(ctx context.Context, pool *pgxpool.Pool, projectID int) ([]Assignee, error) {
	return listProjectAssignees(ctx, pool, projectID)
}

func SetProjectAssignees(ctx context.Context, pool *pgxpool.Pool, projectID int, userIDs []int) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM project_assignees WHERE project_id = $1`, projectID); err != nil {
		return err
	}

	insert := `INSERT INTO project_assignees (project_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	for _, uid := range userIDs {
		if uid <= 0 {
			continue
		}
		if _, err := tx.Exec(ctx, insert, projectID, uid); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func projectOrderBy(sort string, order string) string {
	order = strings.ToLower(strings.TrimSpace(order))
	if order != "asc" {
		order = "desc"
	}

	switch strings.ToLower(strings.TrimSpace(sort)) {
	case "end_date":
		return fmt.Sprintf("p.end_date %s NULLS LAST, p.priority DESC, p.id DESC", order)
	case "start_date":
		return fmt.Sprintf("p.start_date %s NULLS LAST, p.priority DESC, p.id DESC", order)
	case "priority":
		return fmt.Sprintf("p.priority %s, p.end_date ASC NULLS LAST, p.id DESC", order)
	case "created_at":
		return fmt.Sprintf("p.created_at %s, p.id DESC", order)
	case "name":
		return fmt.Sprintf("p.name %s, p.id DESC", order)
	case "urgency", "":
		fallthrough
	default:
		return "p.end_date ASC NULLS LAST, p.priority DESC, p.id DESC"
	}
}

func stageOrderBy(sort string, order string) string {
	order = strings.ToLower(strings.TrimSpace(order))
	if order != "asc" {
		order = "desc"
	}

	switch strings.ToLower(strings.TrimSpace(sort)) {
	case "end_date":
		return fmt.Sprintf("end_date %s NULLS LAST, order_index ASC, id ASC", order)
	case "start_date":
		return fmt.Sprintf("start_date %s NULLS LAST, order_index ASC, id ASC", order)
	case "order_index":
		return fmt.Sprintf("order_index %s, end_date ASC NULLS LAST, id ASC", order)
	case "created_at":
		return fmt.Sprintf("created_at %s, id ASC", order)
	case "urgency", "":
		fallthrough
	default:
		return "end_date ASC NULLS LAST, order_index ASC, id ASC"
	}
}

func progressPercent(done int, total int) int {
	if total <= 0 {
		return 0
	}
	return int(float64(done) / float64(total) * 100)
}
