-- name: DeleteProject :one
UPDATE projects
	SET
		deleted = TRUE,
		username = NULL,
		password = NULL
WHERE id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE)
RETURNING *;
