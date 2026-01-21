-- name: CreateDeployment :one
INSERT INTO deployments (project_id, sha, status, service_entrypoint, branch, env_vars) 
SELECT
	@project_id::bigint AS project_id,
	@sha::text AS sha,
	@status::smallint AS status,
	@service_entrypoint::text AS service_entrypoint,
	@branch::text AS branch,
	@env_vars::jsonb AS env_vars
WHERE NOT EXISTS (SELECT 1
				  FROM deployments
				  WHERE project_id = @project_id::bigint AND branch = @branch::text AND status IN (0, 1) -- New, Deploying
				)
RETURNING *;
