-- name: SelectDeployment :one
SELECT * FROM deployments
WHERE deployments.id = @id::bigint AND project_id = @project_id::bigint;
