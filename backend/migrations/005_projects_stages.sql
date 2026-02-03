CREATE TABLE IF NOT EXISTS projects (
  id SERIAL PRIMARY KEY,
  org_id INTEGER NOT NULL,
  owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  name TEXT NOT NULL,
  description TEXT NULL,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','paused','done')),
  start_date TIMESTAMPTZ NULL,
  end_date TIMESTAMPTZ NULL,
  priority INTEGER NOT NULL DEFAULT 0,
  image_url TEXT NULL,
  budget_allocated NUMERIC(18,2) NOT NULL DEFAULT 0,
  budget_spent NUMERIC(18,2) NOT NULL DEFAULT 0,
  budget_currency TEXT NOT NULL DEFAULT 'KZT',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_projects_org_id ON projects(org_id);
CREATE INDEX IF NOT EXISTS idx_projects_end_priority ON projects(end_date ASC, priority DESC);

CREATE TABLE IF NOT EXISTS project_assignees (
  project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role TEXT NULL,
  PRIMARY KEY (project_id, user_id)
);

CREATE TABLE IF NOT EXISTS stages (
  id SERIAL PRIMARY KEY,
  project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NULL,
  status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo','in_progress','done')),
  start_date TIMESTAMPTZ NULL,
  end_date TIMESTAMPTZ NULL,
  order_index INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stages_project_order ON stages(project_id, order_index ASC);
CREATE INDEX IF NOT EXISTS idx_stages_project_end ON stages(project_id, end_date ASC);
