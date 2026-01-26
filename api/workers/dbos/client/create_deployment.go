package workersdbosclient

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func (c *Client) CreateDeployment(ctx context.Context, params workersinterfaces.CreateDeploymentParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "CreateDeployment", err)

	d := dbosworkflows.Deploy(c.s)

	handle, err := dbos.RunWorkflow(c.ctx, d.Deploy, dbosworkflows.DeployParams{
		ProjectId:    params.ProjectId,
		DeploymentId: params.DeploymentId,
	}, dbos.WithWorkflowID(params.GetWorkflowId()))
	if err != nil {
		return err
	}

	if params.Block {
		_, err = handle.GetResult()
	}
	return err
}
