-- name: UpsertDeploymentLog :one
INSERT INTO deployment_logs (deployment_id, log)
VALUES ((SELECT id FROM deployments WHERE project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE) AND deployments.id = $deploymentId), $log)
ON CONFLICT (deployment_id)
DO UPDATE SET log = CONCAT(COALESCE(deployment_logs.log, ''), CHAR(10), EXCLUDED.log)
RETURNING *;
