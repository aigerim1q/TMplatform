CREATE TABLE IF NOT EXISTS notifications (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- тип уведомления: assigned, deadline_changed, comment, stage_done и т.д.
  type TEXT NOT NULL,

  -- текст/сообщение
  message TEXT NOT NULL,

  -- можно хранить ссылку на сущность (task_id/project_id/stage_id) без жёстких FK
  entity_type TEXT NULL,
  entity_id INTEGER NULL,

  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created
  ON notifications(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_notifications_user_unread
  ON notifications(user_id, is_read);
