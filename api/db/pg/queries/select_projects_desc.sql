-- name: SelectProjectsDesc :many
WITH start AS (
	SELECT MIN(id) AS min
	FROM projects
	ORDER BY id DESC
	LIMIT ($1 * $2) -- limit * page
)

SELECT * FROM projects
WHERE id < (SELECT min FROM start)
ORDER BY id DESC
LIMIT $1;
