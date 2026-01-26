package workersdbosclient

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func (c *Client) BuangBranch(ctx context.Context, params workersinterfaces.BuangBranchParams) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "BuangBranch", err)

	bb := dbosworkflows.BuangBranch(c.s)

	_, err = dbos.RunWorkflow(c.ctx, bb.BuangBranch, dbosworkflows.BuangBranchParams{
		ProjectId: params.ProjectId,
		Branch:    params.Branch,
	})
	return err
}
