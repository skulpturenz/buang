-- name: SelectProjectByRepository :one
SELECT * FROM projects
WHERE repository = @repository::text AND
	  deleted = FALSE;
