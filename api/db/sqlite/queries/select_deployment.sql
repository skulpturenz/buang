-- name: SelectDeployment :one
SELECT * FROM deployments
WHERE deployments.id = $id AND project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE);
