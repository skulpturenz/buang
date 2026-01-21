-- name: UpdateDeploymentStatus :one
UPDATE deployments
	SET status = @status::smallint
WHERE id = @id::bigint AND project_id = @project_id::bigint
RETURNING *;
