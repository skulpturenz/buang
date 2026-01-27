package testutils

import (
	"context"
	"skulpture/buang/db"
	workersshared "skulpture/buang/workers/shared"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
)

func Setup(ctx context.Context, config DurableExecutorConfiguration) (*TestApplication, func(context.Context)) {
	dbCfg := db.DbConfig{
		Type:             config.DbType,
		ConnectionString: config.DbConnectionString,
	}

	queries, dbCleanup, err := dbCfg.New(ctx)
	if err != nil {
		panic(err)
	}

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		dbCleanup(ctx)
		panic(err)
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
	if err != nil {
		docker.Close()
		dbCleanup(ctx)
		panic(err)
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

	cleanup := func(ctx context.Context) {
		appCleanup(ctx)
		workflowsCleanup(ctx)
		docker.Close()
		dbCleanup(ctx)
	}

	return testApp, cleanup
}
