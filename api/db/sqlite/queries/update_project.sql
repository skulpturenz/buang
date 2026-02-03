-- name: UpdateProject :one
UPDATE projects
	SET
		repository = $repository,
		requires_authn = $requiresAuthn,
		username = $username,
		password = $password,
		compose_path = $compose_path
WHERE id = (SELECT id FROM projects WHERE projects.id = $projectId AND projects.deleted = FALSE)
RETURNING *;
