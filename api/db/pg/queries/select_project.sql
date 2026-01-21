-- name: SelectProject :one
SELECT * FROM projects
WHERE id = @id::bigint AND deleted = FALSE;
