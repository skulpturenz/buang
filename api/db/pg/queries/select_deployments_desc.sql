-- name: SelectDeploymentsDesc :many
WITH start AS (SELECT id
		FROM deployments
		WHERE project_id = @project_id::bigint
		ORDER BY id DESC
		-- limit * (page - 1)
		LIMIT (sqlc.arg('limit')::int * (sqlc.arg('page')::int - 1))),
	min AS (SELECT MIN(id) AS min FROM START)


SELECT * FROM deployments
WHERE (deployments.project_id = @project_id::bigint) AND ((sqlc.arg('page')::int = 1) OR (id < (SELECT min FROM min)))
ORDER BY id DESC
LIMIT sqlc.arg('limit')::int;
