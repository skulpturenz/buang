package workers

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
)

type DbosConfig struct {
	Services workersshared.WorkflowServices

	AppName         string
	DatabaseURL     string
	ConductorAPIKey *string
	ConductorURL    *string
	AdminServerPort *int
}

func (c DbosConfig) New(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
	wc := workerConfig{
		Type:     enumsdurableexecutors.Dbos,
		services: c.Services,

		DbosAppName:         c.AppName,
		DbosDatabaseURL:     c.DatabaseURL,
		DbosConductorAPIKey: c.ConductorAPIKey,
		DbosConductorURL:    c.ConductorURL,
		DbosAdminServerPort: c.AdminServerPort,
	}

	return wc.new(ctx)
}
