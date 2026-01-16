CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	deployed_at TIMESTAMP,
	clone_path TEXT,
	CONSTRAINT pk_deployments UNIQUE(id, project_id)
);
