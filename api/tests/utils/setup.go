package testutils

import (
	"context"
	"io"
	"log"
	"skulpture/buang/db"
	"skulpture/buang/ports"
	"skulpture/buang/utils/compensations"
	workersshared "skulpture/buang/workers/shared"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
	"github.com/testcontainers/testcontainers-go"
)

func Setup(ctx context.Context, config DurableExecutorConfiguration, overridePorts ports.Ports) (*TestApplication, func(context.Context)) {
	testcontainers.WithLogger(log.New(io.Discard, "", 0))

	compensations := compensations.New()

	dbCfg := db.DbConfig{
		Type:             config.DbType,
		ConnectionString: config.DbConnectionString,
	}

	queries, dbCleanup, err := dbCfg.New(ctx)
	compensations.AddCompensation(dbCleanup)
	if err != nil {
		compensations.CompensateAndPanic(ctx, err)
	}

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	compensations.AddCompensation(func(ctx context.Context) {
		docker.Close()
	})
	if err != nil {
		compensations.CompensateAndPanic(ctx, err)
	}

	decoder := schema.NewDecoder()
	encoder := schema.NewEncoder()

	var p ports.Ports
	if overridePorts == nil {
		p = ports.New()
	} else {
		p = overridePorts
	}

	ws := workersshared.WorkflowServices{
		Queries:              &queries,
		GorillaSchemaDecoder: decoder,
		GorillaSchemaEncoder: encoder,
		Docker:               docker,
		Ports:                p,
	}
	workflows, workflowsCleanup, err := CreateWorkflows(ctx, config, ws)
	compensations.AddCompensation(workflowsCleanup)
	if err != nil {
		compensations.CompensateAndPanic(ctx, err)
	}

	as := TestApplicationServices{
		Queries:              &queries,
		GorillaSchemaDecoder: decoder,
		GorillaSchemaEncoder: encoder,
		Docker:               docker,
		Workflows:            workflows,
		Ports:                p,
	}

	appConfig := TestApplicationConfig{
		Services:              as,
		DurableExecutorConfig: config,
	}

	testApp, appCleanup := CreateApp(ctx, appConfig)
	compensations.AddCompensation(appCleanup)

	cleanup := compensations.Compensate

	return testApp, cleanup
}
