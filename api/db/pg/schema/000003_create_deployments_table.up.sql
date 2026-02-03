CREATE TABLE IF NOT EXISTS deployments (
	id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
	project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT NOT NULL,
	deployed_at TIMESTAMPTZ,
	clone_path TEXT,
	service_entrypoint TEXT NOT NULL,
	branch TEXT NOT NULL,
	env_vars JSONB,
	CONSTRAINT pk_deployments PRIMARY KEY(id, project_id, branch)
);
