-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint, branch, env_vars) 
SELECT
	$1 AS project_id,
	$2 AS sha,
	$3 AS status,
	$4 AS service_entrypoint,
	$5 AS branch,
	$6 AS env_vars
WHERE NOT EXISTS (SELECT 1
				  FROM deployments
				  WHERE project_id = $1 AND status IN (0, 1) -- New, Deploying
				)
RETURNING *;
