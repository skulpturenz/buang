-- name: DeleteStaleDiagnosticLogs :exec
DELETE FROM diagnostic_logs
WHERE created_at <= DATE('now', '-2 days');
