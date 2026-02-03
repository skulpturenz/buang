package workers

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersdbosclient "skulpture/buang/workers/dbos/client"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
	workerstemporalclient "skulpture/buang/workers/temporal/client"
)

type workerConfig struct {
	Type     enumsdurableexecutors.DurableExecutor
	services workersshared.WorkflowServices

	TemporalHostPort *string

	DbosAppName         string
	DbosDatabaseURL     string
	DbosConductorAPIKey *string
	DbosConductorURL    *string
	DbosAdminServerPort *int
}

func (c workerConfig) new(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
	if c.Type == enumsdurableexecutors.Temporal {
		return workerstemporalclient.New(ctx, workerstemporalclient.TemporalConfig{
			Services: c.services,
			HostPort: c.TemporalHostPort,
		})
	}

	if c.Type == enumsdurableexecutors.Dbos {
		workflows, err := workersdbosclient.New(ctx, workersdbosclient.DbosConfig{
			Services:        c.services,
			AppName:         c.DbosAppName,
			DatabaseURL:     c.DbosDatabaseURL,
			ConductorAPIKey: c.DbosConductorAPIKey,
			ConductorURL:    c.DbosConductorURL,
			AdminServerPort: c.DbosAdminServerPort,
		})
		if err != nil {
			return nil, nil, err
		}

		cleanup := func(ctx context.Context) {} // noop

		return workflows, cleanup, nil
	}

	panic("unsupported durable executor")
}
