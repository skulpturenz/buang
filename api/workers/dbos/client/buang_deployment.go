package workersdbosclient

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func (c *Client) BuangDeployment(ctx context.Context, params workersinterfaces.BuangDeploymentParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "BuangDeployment", err)

	bd := dbosworkflows.BuangDeployment(c.s)

	handle, err := dbos.RunWorkflow(c.ctx, bd.BuangDeployment, dbosworkflows.BuangDeploymentParams{
		ProjectId:    params.ProjectId,
		DeploymentId: params.DeploymentId,
	})
	if err != nil {
		return err
	}

	if params.Block {
		_, err = handle.GetResult()
	}
	return err
}
