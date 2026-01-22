CREATE TABLE IF NOT EXISTS deployment_logs (
	id INTEGER PRIMARY KEY,
	deployment_id INTEGER NOT NULL UNIQUE REFERENCES deployments(id),
	log TEXT,
	CONSTRAINT unique_deployment_logs UNIQUE(id, deployment_id)
);
