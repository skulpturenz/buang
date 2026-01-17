-- name: SelectProjectByRepository :one
SELECT * FROM projects
WHERE repository = $repository AND
	  deleted = 0;
