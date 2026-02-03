-- name: SelectStaleDeployments :many
SELECT * FROM deployments
WHERE deployed_at <= DATE('now', '-2 days');
