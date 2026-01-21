-- name: SelectDeploymentLog :one
SELECT * FROM deployment_logs
WHERE deployment_id = (SELECT id FROM deployments WHERE deployments.id = $deploymentId AND project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE))
