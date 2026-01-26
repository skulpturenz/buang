package workerstemporalclient

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
)

func (c *Client) HasDeployed(ctx context.Context, params workersinterfaces.HasDeployedParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "HasDeployed", err)

	cdp := workersinterfaces.CreateDeploymentParams{
		ProjectId:    params.ProjectId,
		DeploymentId: params.DeploymentId,
	}
	workflowId := cdp.GetWorkflowId()

	run := c.c.GetWorkflow(ctx, workflowId, "")
	err = run.Get(ctx, nil)
	return err
}
