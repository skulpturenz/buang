-- name: CreateProject :one
INSERT INTO projects (repository, requires_authn, username, password, compose_path) 
	VALUES	($repository, $requiresAuthn, $username, $password, $compose_path)
RETURNING *;
