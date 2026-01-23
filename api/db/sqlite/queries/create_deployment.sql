-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint, branch, env_vars) 
SELECT
	(SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE) AS project_id,
	$sha AS sha,
	$status AS status,
	$serviceEntrypoint AS service_entrypoint,
	$branch AS branch,
	$envVars AS env_vars
WHERE NOT EXISTS (SELECT 1
				  FROM deployments
				  WHERE deployments.project_id = (SELECT id FROM projects WHERE projects.id = project_id AND projects.deleted = FALSE) AND deployments.branch = branch AND status IN (0, 1) -- New, Deploying
				  )
RETURNING *;
