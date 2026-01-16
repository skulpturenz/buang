-- name: CreateProject :one
INSERT INTO projects (repository, requires_authn, username, password) 
	VALUES	(?, ?, ?, ?)
RETURNING *;
