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
	ProjectId    int64
	DeploymentId int64
}

func (p CreateDeploymentParams) Exec(ctx context.Context, s app.ApplicationServices) error {
	executor, durableExecutor := s.GetDurableExecutor()

	if durableExecutor == enumsdurableexecutors.Temporal {
		workflowId := fmt.Sprintf("create-project-%v-deployment-%v", p.ProjectId, p.DeploymentId)
		options := temporalclient.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueDeployment,
		}

		_, err := executor.(app.TemporalClient).ExecuteWorkflow(ctx, options, temporalworkflows.Deploy, p.ProjectId, p.DeploymentId)
		if err != nil {
			return err
		}

		return nil
	} else {
		d := dbosworkflows.Deploy(s)

		_, err := dbos.RunWorkflow(executor.(app.DbosContext), d.Deploy, dbosworkflows.DeployParams{
			ProjectId:    p.ProjectId,
			DeploymentId: p.DeploymentId,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
