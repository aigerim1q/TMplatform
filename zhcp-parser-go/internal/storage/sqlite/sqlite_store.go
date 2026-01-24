package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"zhcp-parser-go/internal/storage"
	"zhcp-parser-go/internal/transformers"
)

type Store struct {
	db  *sql.DB
	dsn string
}

func New(dsn string) *Store {
	return &Store{dsn: dsn}
}

func (s *Store) Init(ctx context.Context) error {
	if s.dsn == "" {
		return errors.New("sqlite dsn is required")
	}

	db, err := sql.Open("sqlite", s.dsn)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	// Be conservative: a single writer is usually enough for this prototype.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping sqlite: %w", err)
	}

	s.db = db
	return s.migrate(ctx)
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) PersistProjectStructure(ctx context.Context, structure *transformers.ProjectStructure) (*storage.PersistResult, error) {
	if s.db == nil {
		return nil, errors.New("sqlite store not initialized")
	}
	if structure == nil {
		return nil, errors.New("project structure is nil")
	}

	rawJSON, _ := json.Marshal(structure)

	// Deadline is expected to be on the project; if empty we still store NULL.
	deadline := nullableString(structure.Project.Deadline)

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res := &storage.PersistResult{}

	projectInsert := `
		INSERT INTO projects (title, description, deadline, phases_count, raw_structure_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	pResult, err := tx.ExecContext(ctx, projectInsert,
		structure.Project.Title,
		structure.Project.Description,
		deadline,
		len(structure.Project.Phases),
		string(rawJSON),
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("insert project: %w", err)
	}

	projectID, err := pResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("project last insert id: %w", err)
	}
	res.ProjectID = projectID

	phaseInsert := `
		INSERT INTO phases (project_id, phase_uid, name, description, start_date, end_date, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	taskInsert := `
		INSERT INTO tasks (project_id, phase_id, task_uid, name, description, start_date, end_date, status, responsible_json, dependencies_json, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for i, phase := range structure.Project.Phases {
		pRes, err := tx.ExecContext(ctx, phaseInsert,
			projectID,
			phase.ID,
			phase.Name,
			phase.Description,
			nullableString(phase.StartDate),
			nullableString(phase.EndDate),
			i,
		)
		if err != nil {
			return nil, fmt.Errorf("insert phase %d: %w", i, err)
		}
		phaseID, err := pRes.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("phase last insert id: %w", err)
		}
		res.PhaseIDs = append(res.PhaseIDs, phaseID)

		for j, task := range phase.Tasks {
			responsibleJSON, _ := json.Marshal(task.ResponsiblePersons)
			depsJSON, _ := json.Marshal(task.Dependencies)

			tRes, err := tx.ExecContext(ctx, taskInsert,
				projectID,
				phaseID,
				task.ID,
				task.Name,
				task.Description,
				nullableString(task.StartDate),
				nullableString(task.EndDate),
				task.Status,
				string(responsibleJSON),
				string(depsJSON),
				j,
			)
			if err != nil {
				return nil, fmt.Errorf("insert task %d in phase %d: %w", j, i, err)
			}
			taskID, err := tRes.LastInsertId()
			if err != nil {
				return nil, fmt.Errorf("task last insert id: %w", err)
			}
			res.TaskIDs = append(res.TaskIDs, taskID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return res, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (s *Store) migrate(ctx context.Context) error {
	if s.db == nil {
		return errors.New("db not initialized")
	}

	stmts := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			deadline TEXT NULL,
			phases_count INTEGER NOT NULL,
			raw_structure_json TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS phases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			phase_uid TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			start_date TEXT NULL,
			end_date TEXT NULL,
			order_index INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			phase_id INTEGER NOT NULL,
			task_uid TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			start_date TEXT NULL,
			end_date TEXT NULL,
			status TEXT NOT NULL,
			responsible_json TEXT NOT NULL,
			dependencies_json TEXT NOT NULL,
			order_index INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(phase_id) REFERENCES phases(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_phases_project_id ON phases(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_phase_id ON tasks(phase_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);`,
	}

	for _, stmt := range stmts {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
