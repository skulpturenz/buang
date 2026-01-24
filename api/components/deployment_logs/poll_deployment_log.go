package deploymentlogs

import (
	"context"
	"log/slog"
	"skulpture/buang/app"
	dberrors "skulpture/buang/db/db_errors"
	"skulpture/buang/db/interfaces"
	"time"
)

type PollDeploymentLogParams struct {
	ProjectId    int64
	DeploymentId int64
}

type PollDeploymentLogResult struct {
	Logs string
}

func (p PollDeploymentLogParams) Exec(ctx context.Context, s *app.ApplicationServices) <-chan PollDeploymentLogResult {
	q := *s.Queries

	results := make(chan PollDeploymentLogResult)

	const POLL_INTERVAL = 100 * time.Millisecond

	pollIndefinitely := func() {

		ticker := time.NewTicker(POLL_INTERVAL)
		defer ticker.Stop()
		defer close(results)

		previousLogs := ""
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				res, err := q.SelectDeploymentLog(ctx, interfaces.SelectDeploymentLogParams{
					ProjectID:    p.ProjectId,
					DeploymentID: p.DeploymentId,
				})
				if err != nil && dberrors.IsNoRows(err) {
					continue
				} else if err != nil {
					slog.ErrorContext(ctx, err.Error())
					return
				}

				newLogs := res.GetLog()
				if newLogs != nil && *newLogs != previousLogs {
					previousLogs = *newLogs

					select {
					case results <- PollDeploymentLogResult{Logs: *newLogs}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
	go pollIndefinitely()

	return results
}
