package workflowwrappers

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	temporalclient "go.temporal.io/sdk/client"
)

type BuangDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
	Block        bool
}

func (p BuangDeploymentParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer RecoverWorkflowPanic("BuangDeployment", err)

	executor, durableExecutor := s.GetDurableExecutor()

	if durableExecutor == enumsdurableexecutors.Temporal {
		workflowId := fmt.Sprintf("buang-project-%v-deployment-%v", p.ProjectId, p.DeploymentId)
		options := temporalclient.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueBuang,
		}

		run, err := executor.(app.TemporalClient).ExecuteWorkflow(ctx, options, temporalworkflows.BuangDeployment, p.ProjectId, p.DeploymentId)
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
		bd := dbosworkflows.BuangDeployment(s)

		handle, err := dbos.RunWorkflow(executor.(app.DbosContext), bd.BuangDeployment, dbosworkflows.BuangDeploymentParams{
			ProjectId:    p.ProjectId,
			DeploymentId: p.DeploymentId,
		})
		if err != nil {
			return err
		}

		if p.Block {
			_, err = handle.GetResult()
			if err != nil {
				return err
			}
		}

		return nil
	}
}
