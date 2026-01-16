-- name: CreateProject :one
INSERT INTO projects(repository, requires_authn, username, password) 
	VALUES ($1, $2, $3, $4)
RETURNING *;
