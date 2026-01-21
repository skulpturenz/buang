-- name: SelectDeployment :one
SELECT * FROM deployments
WHERE deployments.id = @id::bigint AND project_id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE);
