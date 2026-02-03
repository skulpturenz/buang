CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY,
	repository TEXT NOT NULL,
	requires_authn BOOLEAN NOT NULL DEFAULT 0,
	username TEXT,
	password TEXT,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	compose_path TEXT NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT 0,
	CONSTRAINT unique_projects UNIQUE(id, repository)
);

CREATE TRIGGER IF NOT EXISTS projects_moddatetime
	BEFORE UPDATE ON projects
	FOR EACH ROW BEGIN
		UPDATE projects
			SET updated_at = CURRENT_TIMESTAMP 
		WHERE rowid = OLD.rowid;
	END;
