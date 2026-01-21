-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = @status::smallint
WHERE id = @id::bigint AND project_id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE)
RETURNING *;
