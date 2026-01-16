CREATE TABLE IF NOT EXISTS deployments (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	project_id BIGINT NOT NULL REFERENCES projects(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	deployed_at TIMESTAMPTZ,
	CONSTRAINT pk_deployments PRIMARY KEY(id, project_id)
);
