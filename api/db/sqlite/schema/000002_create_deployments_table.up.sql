CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER NOT NULL AUTOINCREMENT,
	repository_id INTEGER NOT NULL REFERENCES repositories(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	CONSTRAINT pk_deployments PRIMARY KEY(id, repository_id)
);
