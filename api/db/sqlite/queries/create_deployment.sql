-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint, branch, env_vars) 
	VALUES	($projectId, $sha, $status, $serviceEntrypoint, $branch, $envVars)
RETURNING *;
