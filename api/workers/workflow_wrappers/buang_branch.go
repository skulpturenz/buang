package workflowwrappers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	temporal "go.temporal.io/sdk/client"
)

type BuangBranchParams struct {
	ProjectId int64
	Branch    string
}

func (p BuangBranchParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang branch panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("buang branch panic")
			}
		}
	}()

	executor, durableExecutor := s.GetDurableExecutor()

	if durableExecutor == enumsdurableexecutors.Temporal {
		workflowId := fmt.Sprintf("buang-project-%v-branch-%v", p.ProjectId, p.Branch)
		options := temporal.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueBuang,
		}

		_, err := executor.(app.TemporalClient).ExecuteWorkflow(ctx, options, temporalworkflows.BuangBranch, int64(p.ProjectId), p.Branch)
		if err != nil {
			return err
		}

		return nil
	} else {
		bb := dbosworkflows.BuangBranch(s)

		_, err := dbos.RunWorkflow(executor.(dbos.DBOSContext), bb.BuangBranch, dbosworkflows.BuangBranchParams{
			ProjectId: p.ProjectId,
			Branch:    p.Branch,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
