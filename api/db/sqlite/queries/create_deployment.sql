-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint) 
	VALUES	($projectId, $sha, $status, $service_entrypoint)
RETURNING *;
