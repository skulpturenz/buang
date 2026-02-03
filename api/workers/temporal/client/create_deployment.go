package workerstemporalclient

import (
	"context"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	temporalclient "go.temporal.io/sdk/client"
)

func (c *Client) CreateDeployment(ctx context.Context, params workersinterfaces.CreateDeploymentParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "CreateDeployment", err)

	options := temporalclient.StartWorkflowOptions{
		ID:        params.GetWorkflowId(),
		TaskQueue: constantstaskqueues.QueueDeployment,
	}

	run, err := c.c.ExecuteWorkflow(ctx, options, temporalworkflows.Deploy, params.ProjectId, params.DeploymentId)
	if err != nil {
		return err
	}

	if params.Block {
		err = run.Get(ctx, nil)
	}
	return err
}
