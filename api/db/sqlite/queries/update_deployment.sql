-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = $url,
		status = $status,
		deployed_at = $deployedAt,
		clone_path = $clonePath
WHERE deployments.id = $id AND project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)
RETURNING *;
