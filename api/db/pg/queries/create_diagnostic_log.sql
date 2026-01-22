-- name: CreateDiagnosticLog :one
INSERT INTO diagnostic_logs(type, log) 
	VALUES (@type::text, @log::jsonb)
RETURNING *;
