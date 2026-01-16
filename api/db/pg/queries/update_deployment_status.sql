-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = $2
WHERE id = $1
RETURNING *;
