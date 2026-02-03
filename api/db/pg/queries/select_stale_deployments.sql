-- name: SelectStaleDeployments :many
SELECT * FROM deployments
WHERE deployed_at <= NOW() - INTERVAL '2 days';
