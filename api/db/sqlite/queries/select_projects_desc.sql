-- name: SelectProjectsDesc :many
WITH start AS (
	SELECT MIN(id) AS min
	FROM projects
	ORDER BY id DESC
	LIMIT ($limit * $page) -- limit * page
)

SELECT * FROM projects
WHERE id < (SELECT min FROM start)
ORDER BY id DESC
LIMIT $limit
