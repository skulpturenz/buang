package deploymentlogs

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	"time"

	"github.com/jackc/pgx/v5"
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
				if err != nil && (errors.Is(sql.ErrNoRows, err) || errors.Is(pgx.ErrNoRows, err)) {
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
