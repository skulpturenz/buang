package testutils

import (
	"context"
	"io"
	"log"
	"skulpture/buang/db"
	"skulpture/buang/utils/compensations"
	workersshared "skulpture/buang/workers/shared"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
	"github.com/testcontainers/testcontainers-go"
)

func Setup(ctx context.Context, config DurableExecutorConfiguration) (*TestApplication, func(context.Context)) {
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

	ws := workersshared.WorkflowServices{
		Queries:              &queries,
		GorillaSchemaDecoder: decoder,
		GorillaSchemaEncoder: encoder,
		Docker:               docker,
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
