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

func (c *Client) BuangBranch(ctx context.Context, params workersinterfaces.BuangBranchParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "BuangBranch", err)

	workflowId := fmt.Sprintf("buang-project-%v-branch-%v", params.ProjectId, params.Branch)
	options := temporalclient.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: constantstaskqueues.QueueBuang,
	}

	_, err = c.c.ExecuteWorkflow(ctx, options, temporalworkflows.BuangBranch, int64(params.ProjectId), params.Branch)
	return err
}
