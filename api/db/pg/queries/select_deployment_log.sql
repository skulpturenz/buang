-- name: SelectDeploymentLog :one
SELECT * FROM deployment_logs
WHERE deployment_id = (SELECT id FROM deployments WHERE id = @deployment_id::bigint AND project_id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE));
