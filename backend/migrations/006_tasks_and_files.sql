CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  stage_id INTEGER NOT NULL REFERENCES stages(id) ON DELETE CASCADE,
  parent_id INTEGER NULL REFERENCES tasks(id) ON DELETE CASCADE, -- For subtasks
  assignee_id INTEGER NULL REFERENCES users(id) ON DELETE SET NULL,
  author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  title TEXT NOT NULL,
  description TEXT NULL,
  status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo','in_progress','done','review')),
  priority INTEGER NOT NULL DEFAULT 0, -- 0: Low, 1: Medium, 2: High, 3: Updates
  start_date TIMESTAMPTZ NULL,
  due_date TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tasks_stage_id ON tasks(stage_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks(assignee_id);

CREATE TABLE IF NOT EXISTS files (
  id SERIAL PRIMARY KEY,
  task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  uploader_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  filename TEXT NOT NULL,       -- Original name "report.pdf"
  storage_path TEXT NOT NULL,   -- Path on disk "uploads/uuid-report.pdf" or S3 key
  mime_type TEXT NOT NULL,
  size_bytes BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_files_task_id ON files(task_id);
