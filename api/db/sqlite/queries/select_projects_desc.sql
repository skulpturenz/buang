-- name: SelectProjectsDesc :many
WITH start AS (SELECT id
		FROM projects
		ORDER BY id DESC
		-- limit * (page - 1)
		LIMIT ($limit * ($page - 1))),
	min AS (SELECT MIN(id) AS min FROM START)


SELECT * FROM projects
WHERE ($page = 1) OR (id < (SELECT min FROM min))
ORDER BY id DESC
LIMIT $limit
