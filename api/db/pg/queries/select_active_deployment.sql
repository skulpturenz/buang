-- name: SelectActiveDeployments :many
SELECT * FROM deployments
WHERE project_id = $1 AND 
	  status IN (0, 1, 2); -- New, Deploying, Deployed
