package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Migrate applies idempotent schema creation for the chatbot project database.
func Migrate(dbConn *sql.DB) error {
	stmts := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT UNIQUE,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            CHECK (length(trim(name)) > 0)
        );`,
		`CREATE TABLE IF NOT EXISTS projects (
		    id TEXT PRIMARY KEY,
		    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		    title TEXT NOT NULL,
		    description TEXT NOT NULL DEFAULT '',
		    status TEXT NOT NULL DEFAULT 'active',
		    normalized_title TEXT NOT NULL,
		    next_task_stable_id INTEGER NOT NULL DEFAULT 1,  -- Counter for next stable task ID
		    deleted_at TIMESTAMP,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    CHECK (length(trim(title)) > 0),
		    CHECK (status IN ('active','inactive','archived','closed')),
		    UNIQUE(owner_id, normalized_title)
		);`,
		`CREATE TABLE IF NOT EXISTS project_members (
            project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
            user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            role TEXT NOT NULL DEFAULT 'member',
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY(project_id, user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS stages (
            id TEXT PRIMARY KEY,
            project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
            title TEXT NOT NULL,
            normalized_title TEXT NOT NULL,
            order_index INTEGER NOT NULL DEFAULT 0,
            responsible_id TEXT REFERENCES users(id) ON DELETE SET NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            CHECK (length(trim(title)) > 0),
            UNIQUE(project_id, normalized_title)
        );`,
		`CREATE TABLE IF NOT EXISTS tasks (
		    id TEXT PRIMARY KEY,
		    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		    stage_id TEXT REFERENCES stages(id) ON DELETE SET NULL,
		    title TEXT NOT NULL,
		    normalized_title TEXT NOT NULL,
		    description TEXT NOT NULL DEFAULT '',
		    status TEXT NOT NULL DEFAULT 'open',
		    priority TEXT NOT NULL DEFAULT 'medium',
		    assignee_id TEXT REFERENCES users(id) ON DELETE SET NULL,
		    due_date TIMESTAMP,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    numeric_id INTEGER,  -- Stable numeric ID for referencing tasks
		    CHECK (length(trim(title)) > 0),
		    CHECK (status IN ('open','in_progress','done','blocked')),
		    CHECK (priority IN ('low','medium','high')),
		    UNIQUE(project_id, normalized_title)
		);`,
		`CREATE TABLE IF NOT EXISTS idempotency_keys (
            key TEXT PRIMARY KEY,
            user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            intent TEXT NOT NULL,
            resource_type TEXT NOT NULL,
            resource_id TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE INDEX IF NOT EXISTS idx_projects_owner ON projects(owner_id);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_normalized_title ON projects(normalized_title);`,
		`CREATE INDEX IF NOT EXISTS idx_stages_project ON stages(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_stage ON tasks(stage_id);`,
		`CREATE INDEX IF NOT EXISTS idx_idempotency_user ON idempotency_keys(user_id);`,
	}

	for _, stmt := range stmts {
		if _, err := dbConn.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Add new columns in an idempotent way for existing installations.
	if _, err := dbConn.Exec(`ALTER TABLE projects ADD COLUMN deleted_at TIMESTAMP`); err != nil {
		// SQLite returns an error if the column already exists; ignore that specific case.
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Add numeric_id column to tasks table if it doesn't exist
	if _, err := dbConn.Exec(`ALTER TABLE tasks ADD COLUMN numeric_id INTEGER`); err != nil {
		// SQLite returns an error if the column already exists; ignore that specific case.
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Add next_task_stable_id column to projects table if it doesn't exist
	if _, err := dbConn.Exec(`ALTER TABLE projects ADD COLUMN next_task_stable_id INTEGER NOT NULL DEFAULT 1`); err != nil {
		// SQLite returns an error if the column already exists; ignore that specific case.
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
