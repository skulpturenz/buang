-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = $2,
		status = $3,
		deployed_at = $4
WHERE id = $1
RETURNING *;
