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

type CreateDeploymentParams struct {
	ProjectId    int64 // required for workflow id
	DeploymentId int64 // required for workflow id
	Block        bool
}

func (p CreateDeploymentParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer RecoverWorkflowPanic("CreateDeployment", err)

	executor, durableExecutor := s.GetDurableExecutor()

	if durableExecutor == enumsdurableexecutors.Temporal {
		options := temporalclient.StartWorkflowOptions{
			ID:        p.GetWorkflowId(),
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
		}, dbos.WithWorkflowID(p.GetWorkflowId()))
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

func (p CreateDeploymentParams) GetWorkflowId() string {
	return fmt.Sprintf("create-project-%v-deployment-%v", p.ProjectId, p.DeploymentId)
}
