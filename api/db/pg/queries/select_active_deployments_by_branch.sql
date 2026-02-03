-- name: SelectActiveDeploymentsByBranch :many
SELECT * FROM deployments
WHERE project_id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE) AND 
	  branch = @branch::text AND
	  clone_path IS NOT NULL AND
	  status IN (0, 1, 2); -- New, Deploying, Deployed
