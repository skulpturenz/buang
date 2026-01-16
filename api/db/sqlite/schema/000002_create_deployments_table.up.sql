CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER PRIMARY KEY,
	repository_id INTEGER NOT NULL REFERENCES repositories(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	deployed_at TIMESTAMP,
	CONSTRAINT pk_deployments UNIQUE(id, repository_id)
);
