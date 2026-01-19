CREATE TABLE IF NOT EXISTS deployment_logs (
	id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
	deployment_id BIGINT NOT NULL REFERENCES deployments(id),
	log TEXT,
	CONSTRAINT pk_deployment_logs PRIMARY KEY(id, deployment_id)
);
