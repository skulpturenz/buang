package dbos

import (
	"skulpture/buang/app"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Deployment(s app.ApplicationServices, ctx dbos.DBOSContext) {
	dbos.RegisterWorkflow(ctx, dbosworkflows.Deploy)
}
