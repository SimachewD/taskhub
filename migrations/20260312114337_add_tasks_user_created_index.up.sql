CREATE INDEX idx_tasks_user_created
ON tasks(user_id, created_at DESC, id DESC);