-- name: SelectActiveDeploymentsByBranch :many
SELECT * FROM deployments
WHERE project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE) AND
	  branch = $branch AND
	  clone_path IS NOT NULL AND
	  status IN (0, 1, 2); -- New, Deploying, Deployed
