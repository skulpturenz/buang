CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id),
	url TEXT,
	status SMALLINT NOT NULL,
	sha TEXT NOT NULL,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	deployed_at TIMESTAMP,
	clone_path TEXT,
	service_entrypoint TEXT NOT NULL,
	branch TEXT NOT NULL,
	env_vars BLOB,
	CONSTRAINT pk_deployments UNIQUE(id, project_id, branch)
);
