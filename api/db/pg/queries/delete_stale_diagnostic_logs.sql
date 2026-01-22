-- name: DeleteStaleDiagnosticLogs :exec
DELETE FROM diagnostic_logs
WHERE created_at <= NOW() - INTERVAL '2 days';
