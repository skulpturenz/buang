-- name: DeleteStaleDiagnosticLogs :exec
DELETE FROM diagnostic_logs
WHERE created_at <= DATE('now', '-2 days') AND
	  1000 <= (SELECT COUNT(1) FROM diagnostic_logs);
