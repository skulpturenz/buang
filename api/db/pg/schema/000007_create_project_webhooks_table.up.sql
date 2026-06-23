CREATE TABLE IF NOT EXISTS project_webhooks (
	id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
	webhook_id BIGINT NOT NULL,
	webhook_type SMALLINT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	project_id BIGINT NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT pk_project_webhooks PRIMARY KEY(id, webhook_id, webhook_type, project_id),
	CONSTRAINT fk_project_webhooks_webhook FOREIGN KEY (webhook_id, webhook_type) REFERENCES webhooks(id, type),
	CONSTRAINT fk_project_webhooks_project FOREIGN KEY (project_id) REFERENCES projects(id)
);