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
	temporalclient "go.temporal.io/sdk/client"
)

type CreateDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
	Block        bool
}

func (p CreateDeploymentParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("create deployment panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("create deployment panic")
			}
		}
	}()

	executor, durableExecutor := s.GetDurableExecutor()

	if durableExecutor == enumsdurableexecutors.Temporal {
		workflowId := fmt.Sprintf("create-project-%v-deployment-%v", p.ProjectId, p.DeploymentId)
		options := temporalclient.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueDeployment,
		}

		run, err := executor.(app.TemporalClient).ExecuteWorkflow(ctx, options, temporalworkflows.Deploy, p.ProjectId, p.DeploymentId)
		if err != nil {
			return err
		}

		if p.Block {
			err = run.Get(ctx, nil)
			if err != nil {
				return err
			}
		}

		return nil
	} else {
		d := dbosworkflows.Deploy(s)

		handle, err := dbos.RunWorkflow(executor.(app.DbosContext), d.Deploy, dbosworkflows.DeployParams{
			ProjectId:    p.ProjectId,
			DeploymentId: p.DeploymentId,
		})
		if err != nil {
			return err
		}

		if p.Block {
			_, err := handle.GetResult()
			if err != nil {
				return err
			}
		}

		return nil
	}
}
