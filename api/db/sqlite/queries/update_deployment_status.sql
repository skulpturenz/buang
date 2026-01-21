-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = $status
WHERE id = $id AND project_id = $projectId
RETURNING *;
