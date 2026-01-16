-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status) 
	VALUES	($1, $2, $3)
RETURNING *;
