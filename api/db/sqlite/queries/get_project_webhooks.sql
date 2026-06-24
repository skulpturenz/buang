-- name: GetProjectWebhooks :many
SELECT pw.id, pw.webhook_id, pw.webhook_type, pw.url, pw.project_id, pw.deleted
FROM project_webhooks pw
JOIN webhooks w ON pw.webhook_id = w.id AND pw.webhook_type = w.type
WHERE pw.project_id = $projectId AND pw.deleted = 0 AND w.deleted = 0;