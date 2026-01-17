-- name: CreateProject :one
INSERT INTO projects(repository, requires_authn, username, password, compose_path) 
	VALUES ($1, $2, $3, $4, $5)
RETURNING *;
