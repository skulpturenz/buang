-- name: SelectProject :one
SELECT * FROM projects
WHERE id = $id AND deleted = FALSE;
