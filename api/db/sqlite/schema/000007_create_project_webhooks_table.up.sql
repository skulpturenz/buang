CREATE TABLE IF NOT EXISTS project_webhooks (
	id INTEGER PRIMARY KEY,
	webhook_id INTEGER NOT NULL,
	webhook_type INTEGER NOT NULL,
	url TEXT NOT NULL UNIQUE,
	project_id INTEGER NOT NULL,
	deleted INTEGER NOT NULL DEFAULT 0,
	CONSTRAINT unique_project_webhooks UNIQUE(id, webhook_id, webhook_type, project_id),
	CONSTRAINT fk_project_webhooks_webhook FOREIGN KEY (webhook_id, webhook_type) REFERENCES webhooks(id, type),
	CONSTRAINT fk_project_webhooks_project FOREIGN KEY (project_id) REFERENCES projects(id)
);