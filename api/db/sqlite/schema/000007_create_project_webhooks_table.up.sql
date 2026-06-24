CREATE TABLE IF NOT EXISTS project_webhooks (
	id INTEGER PRIMARY KEY,
	webhook_id INTEGER NOT NULL REFERENCES webhooks(id),
	url TEXT NOT NULL UNIQUE,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	deleted INTEGER NOT NULL DEFAULT 0,
	CONSTRAINT unique_project_webhooks UNIQUE(id, webhook_id, project_id),
);
