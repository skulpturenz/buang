package workerstemporalclient

import (
	"context"
	"fmt"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	temporalclient "go.temporal.io/sdk/client"
)

func (c *Client) BuangDeployment(ctx context.Context, params workersinterfaces.BuangDeploymentParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "BuangDeployment", err)

	workflowId := fmt.Sprintf("buang-project-%v-deployment-%v", params.ProjectId, params.DeploymentId)
	options := temporalclient.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: constantstaskqueues.QueueBuang,
	}

	run, err := c.c.ExecuteWorkflow(ctx, options, temporalworkflows.BuangDeployment, params.ProjectId, params.DeploymentId)
	if err != nil {
		return err
	}

	if params.Block {
		err = run.Get(ctx, nil)
	}
	return err
}
