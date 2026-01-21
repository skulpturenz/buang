-- name: SelectDeploymentsDesc :many
WITH start AS (SELECT id
		FROM deployments
		WHERE project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)
		ORDER BY id DESC
		-- limit * (page - 1)
		LIMIT ($limit * ($page - 1))),
	min AS (SELECT MIN(id) AS min FROM START)


SELECT * FROM deployments
WHERE (deployments.project_id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)) AND (($page = 1) OR (id < (SELECT min FROM min)))
ORDER BY id DESC
LIMIT $limit
