package dbos

import (
	"skulpture/buang/app"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Housekeping(s app.ApplicationServices, ctx dbos.DBOSContext) {
	dbos.RegisterWorkflow(ctx, dbosworkflows.BuangHousekeeping,
		dbos.WithSchedule("0 0 */2 * *")) // every 2 days
}
