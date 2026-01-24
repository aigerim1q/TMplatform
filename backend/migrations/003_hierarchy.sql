DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'users'
      AND column_name = 'manager_id'
      AND data_type = 'uuid'
  ) THEN
    ALTER TABLE users
      ALTER COLUMN manager_id TYPE INTEGER USING NULL;
  END IF;
END $$;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS manager_id INTEGER NULL;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_manager') THEN
    ALTER TABLE users
      ADD CONSTRAINT fk_users_manager
      FOREIGN KEY (manager_id) REFERENCES users(id)
      ON DELETE SET NULL;
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_manager_not_self') THEN
    ALTER TABLE users
      ADD CONSTRAINT chk_manager_not_self
      CHECK (manager_id IS NULL OR manager_id <> id);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_manager_id ON users(manager_id);
