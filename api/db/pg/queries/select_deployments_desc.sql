-- name: SelectDeploymentsDesc :many
WITH start AS (SELECT id
		FROM deployments
		WHERE project_id = $1
		ORDER BY id DESC
		-- limit * (page - 1)
		LIMIT ($2 * ($3 - 1))),
	min AS (SELECT MIN(id) AS min FROM START)


SELECT * FROM deployments
WHERE (deployments.project_id = $1) AND (($3 = 1) OR (id < (SELECT min FROM min)))
ORDER BY id DESC
LIMIT $2;
