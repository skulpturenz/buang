package workersdbos

import (
	"skulpture/buang/app"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Deployment(s app.ApplicationServices, ctx dbos.DBOSContext) {
	d := dbosworkflows.Deploy(s)
	dbos.RegisterWorkflow(ctx, d.Deploy, dbos.WithMaxRetries(10))
}
