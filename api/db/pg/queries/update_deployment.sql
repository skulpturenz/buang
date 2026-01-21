-- name: UpdateDeployment :one
UPDATE deployments
	SET 
		url = sqlc.narg('url'),
		status = @status::smallint,
		deployed_at = sqlc.narg('deployed_at'),
		clone_path = sqlc.narg('clone_path')
WHERE id = @id::bigint AND project_id = (SELECT id FROM projects WHERE projects.id = @project_id::bigint)
RETURNING *;
