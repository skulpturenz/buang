-- name: UpsertDeploymentLog :one
INSERT INTO deployment_logs (deployment_id, log)
VALUES ((SELECT id FROM deployments WHERE project_id = $1 AND deployments.id = $2), $3)
ON CONFLICT (deployment_id)
DO UPDATE SET log = CONCAT(COALESCE(deployment_logs.log, ''), CHR(10), EXCLUDED.log)
RETURNING *;
