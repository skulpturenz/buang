-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = $status
WHERE deployments.id = $id AND project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)
RETURNING *;
