-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status) 
	VALUES	($projectId, $sha, $status)
RETURNING *;
