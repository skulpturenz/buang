CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER PRIMARY KEY,
	repository_id INTEGER NOT NULL REFERENCES repositories(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	CONSTRAINT pk_deployments UNIQUE(id, repository_id)
);
