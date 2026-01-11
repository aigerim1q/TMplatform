ALTER TABLE users
  ADD COLUMN IF NOT EXISTS manager_id uuid NULL;

ALTER TABLE users
  ADD CONSTRAINT IF NOT EXISTS fk_users_manager
  FOREIGN KEY (manager_id) REFERENCES users(id)
  ON DELETE SET NULL;

ALTER TABLE users
  ADD CONSTRAINT IF NOT EXISTS chk_manager_not_self
  CHECK (manager_id IS NULL OR manager_id <> id);

CREATE INDEX IF NOT EXISTS idx_users_manager_id ON users(manager_id);
