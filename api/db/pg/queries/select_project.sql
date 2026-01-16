-- name: SelectProject :one
SELECT * FROM projects
WHERE id = $1 AND deleted = FALSE;
