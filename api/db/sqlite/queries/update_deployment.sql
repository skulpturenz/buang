-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = $url,
		status = $status,
		deployed_at = $deployedAt,
		clone_path = $clonePath
WHERE id = $id
RETURNING *;
