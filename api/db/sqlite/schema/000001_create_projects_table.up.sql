CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY,
	repository TEXT NOT NULL,
	requires_authn BOOLEAN NOT NULL DEFAULT 0,
	username TEXT,
	password TEXT,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT unique_projects UNIQUE(id, repository)
);

CREATE TRIGGER IF NOT EXISTS projects_moddatetime
	BEFORE UPDATE ON projects
	FOR EACH ROW BEGIN
		UPDATE projects
			SET updated_at = CURRENT_TIMESTAMP 
		WHERE rowid = OLD.rowid;
	END;
