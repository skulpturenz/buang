-- name: UpdateProject :one
UPDATE projects
	SET
		repository = @repository::text,
		requires_authn = @requires_authn::bool,
		username = sqlc.narg('username')::text,
		password = sqlc.narg('password')::text,
		compose_path = @compose_path::text
WHERE id = (SELECT id FROM projects WHERE id = @project_id::bigint AND deleted = FALSE)
RETURNING *;
