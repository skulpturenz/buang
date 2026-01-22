package dbos

import (
	"os"
	"skulpture/buang/app"
	constantsfeaturetoggles "skulpture/buang/constants/feature_toggles"
	enumsenv "skulpture/buang/enums/env"
	dbosworkflows "skulpture/buang/workers/dbos/workflows"
	"strconv"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

func Housekeping(s app.ApplicationServices, ctx dbos.DBOSContext) {
	bh := dbosworkflows.BuangHousekeeping(s)

	dbos.RegisterWorkflow(ctx, bh.BuangHousekeeping,
		dbos.WithSchedule("0 0 */2 * * *"), // every 2 days
		dbos.WithMaxRetries(3))

	goEnv, _ := os.LookupEnv("GO_ENV")
	goEnvE, _ := enumsenv.Parse(goEnv)
	isExperimentalBootstrapEnabledEnv, _ := os.LookupEnv(constantsfeaturetoggles.EXPERIMENTAL_BOOTSTRAP)
	isExperimentalBootstrapEnabled, _ := strconv.ParseBool(isExperimentalBootstrapEnabledEnv)

	if goEnvE == enumsenv.Production && isExperimentalBootstrapEnabled {
		ab := dbosworkflows.PeriodicUpdateHandler(s)

		dbos.RegisterWorkflow(ctx, ab.PeriodicUpdateHandler,
			dbos.WithSchedule("0 0 * * 6 *"), // every saturday
			dbos.WithMaxRetries(3))
	}
}
