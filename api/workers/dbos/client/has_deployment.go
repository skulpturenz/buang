package workersdbosclient

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func (c *Client) HasDeployed(ctx context.Context, params workersinterfaces.HasDeployedParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "HasDeployed", err)

	cdp := workersinterfaces.CreateDeploymentParams{
		ProjectId:    params.ProjectId,
		DeploymentId: params.DeploymentId,
	}

	handle, err := dbos.RetrieveWorkflow[bool](c.ctx, cdp.GetWorkflowId())
	if err != nil {
		return err
	}

	_, err = handle.GetResult()
	return err
}
