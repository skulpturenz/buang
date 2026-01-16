-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = $url,
		status = $status,
		deployed_at = $deployedAt
WHERE id = $id
RETURNING *;
