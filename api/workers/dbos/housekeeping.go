package dbos

import (
	"skulpture/buang/app"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Housekeping(s app.ApplicationServices, ctx dbos.DBOSContext) {
	bh := dbosworkflows.BuangHousekeeping(s)

	dbos.RegisterWorkflow(ctx, bh.BuangHousekeeping,
		dbos.WithSchedule("0 0 */2 * * *"), // every 2 days
		dbos.WithMaxRetries(3))
}
