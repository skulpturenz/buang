-- name: SelectDeployment :one
SELECT * FROM deployments
WHERE id = $id AND project_id = $projectId;
