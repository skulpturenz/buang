-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = $2,
		status = $3,
		deployed_at = $4,
		clone_path = $5
WHERE id = $1
RETURNING *;
