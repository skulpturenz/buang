-- name: UpsertDeploymentLog :one
INSERT INTO deployment_logs (deployment_id, log)
VALUES ((SELECT id FROM deployments WHERE project_id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE) AND deployments.id = @id::bigint), sqlc.narg('log'))
ON CONFLICT (deployment_id)
DO UPDATE SET log = CONCAT(COALESCE(deployment_logs.log, ''), CHR(10), EXCLUDED.log)
RETURNING *;
