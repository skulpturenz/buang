package dbos

import (
	"skulpture/buang/app"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Buang(s app.ApplicationServices, ctx dbos.DBOSContext) {
	bd := dbosworkflows.BuangDeployment(s)
	dbos.RegisterWorkflow(ctx, bd.BuangDeployment)

	bb := dbosworkflows.BuangBranch(s)
	dbos.RegisterWorkflow(ctx, bb.BuangBranch)
}
