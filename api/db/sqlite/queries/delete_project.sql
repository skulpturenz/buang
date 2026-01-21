-- name: DeleteProject :one
UPDATE projects
	SET
		deleted = TRUE,
		username = NULL,
		password = NULL
WHERE id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)
RETURNING *;
