-- name: CreateDiagnosticLog :one
INSERT INTO diagnostic_logs(type, log) 
	VALUES ($type, $log)
RETURNING *;
