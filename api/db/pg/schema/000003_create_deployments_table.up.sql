CREATE TABLE IF NOT EXISTS deployments (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	repository_id BIGINT NOT NULL REFERENCES repositories(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	deployed_at TIMESTAMPTZ,
	CONSTRAINT pk_deployments PRIMARY KEY(id, repository_id)
);
