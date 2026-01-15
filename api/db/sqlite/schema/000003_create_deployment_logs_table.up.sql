CREATE TABLE IF NOT EXISTS deployment_logs (
	id INTEGER NOT NULL AUTOINCREMENT,
	deployment_id INTEGER NOT NULL REFERENCES deployments(id),
	log TEXT,
	CONSTRAINT pk_deployment_logs PRIMARY KEY(id, deployment_id)
);
