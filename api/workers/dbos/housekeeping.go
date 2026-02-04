package workersdbos

import (
	"fmt"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	enumsenv "skulpture/buang/enums/env"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Housekeping(s app.ApplicationServices, ctx dbos.DBOSContext) {
	bh := dbosworkflows.BuangHousekeeping(s)

	dbos.RegisterWorkflow(ctx, bh.BuangHousekeeping,
		dbos.WithSchedule(fmt.Sprintf("%v *", constantsenvs.HOUSEKEEPING_PRUNE_DEPLOYMENTS.Value())),
		dbos.WithMaxRetries(10))

	goEnv, _ := enumsenv.Parse(constantsenvs.GO_ENV.Value())
	isBootstrapEnabled, _ := constantsenvs.EXPERIMENTAL_BOOTSTRAP.Value()

	if goEnv == enumsenv.Production && isBootstrapEnabled {
		ab := dbosworkflows.PeriodicUpdateHandler(s)

		schedule, _ := constantsenvs.EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE.Value()

		dbos.RegisterWorkflow(ctx, ab.PeriodicUpdateHandler,
			dbos.WithSchedule(fmt.Sprintf("%v *", schedule)),
			dbos.WithMaxRetries(10))
	}
}
