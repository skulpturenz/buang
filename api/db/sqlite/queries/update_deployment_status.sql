-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = $status
WHERE id = $id
RETURNING *;
