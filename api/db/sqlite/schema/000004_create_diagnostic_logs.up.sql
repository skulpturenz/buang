CREATE TABLE IF NOT EXISTS diagnostic_logs (
	id INTEGER PRIMARY KEY,
	type TEXT NOT NULL,
	-- https://github.com/mattn/go-sqlite3/blob/master/doc.go#L24
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	log BLOB NOT NULL,
	CONSTRAINT unique_diagnostic_logs UNIQUE(id, type)
);
