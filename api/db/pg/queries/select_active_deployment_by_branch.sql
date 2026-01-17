-- name: SelectActiveDeploymentsByBranch :many
SELECT * FROM deployments
WHERE project_id = $1 AND 
	  branch = $2 AND
	  clone_path IS NOT NULL AND
	  status IN (0, 1, 2); -- New, Deploying, Deployed
