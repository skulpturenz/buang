-- name: SelectActiveDeploymentsByProject :many
SELECT * FROM deployments
WHERE project_id = (SELECT id FROM projects WHERE projects.id = @projectId::bigint AND projects.deleted = FALSE) AND
	  status IN (0, 1, 2); -- New, Deploying, Deployed
