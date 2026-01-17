-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint) 
	VALUES	($1, $2, $3, $4)
RETURNING *;
