-- name: SelectProjectsDesc :many
WITH start AS (
	SELECT MIN(id) AS min
	FROM projects
	WHERE deleted = FALSE
	ORDER BY id DESC
	-- limit * (page - 1)
	LIMIT ($1 * ($2 - 1))
)

SELECT * FROM projects
WHERE (deleted = FALSE) AND (($2 = 1) OR (id < (SELECT min FROM start)))
ORDER BY id DESC
LIMIT $1;
