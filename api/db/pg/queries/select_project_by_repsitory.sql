-- name: SelectProjectByRepository :one
SELECT * FROM projects
WHERE repository = $1 AND
	  deleted = 0;
