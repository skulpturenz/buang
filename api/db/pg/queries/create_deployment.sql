-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint, branch, env_vars) 
	VALUES	($1, $2, $3, $4, $5, $6)
RETURNING *;
