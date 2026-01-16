-- name: SelectDeployment :one
SELECT * FROM deployments
WHERE id = $1 AND project_id = $2;
