package testutils

import (
	"context"
	"io"
	"log"
	"skulpture/buang/db"
	"skulpture/buang/ports"
	"skulpture/buang/services"
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

	svc := services.New()
	services.Set(svc, services.KeyQueries, queries)
	services.Set(svc, services.KeyDocker, docker)
	services.Set(svc, services.KeyGorillaSchemaDecoder, *decoder)
	services.Set(svc, services.KeyGorillaSchemaEncoder, *encoder)
	services.Set(svc, services.KeyPorts, p)

	ws := workersshared.WorkflowServices{
		Services: svc,
	}
	workflows, workflowsCleanup, err := CreateWorkflows(ctx, config, ws)
	compensations.AddCompensation(workflowsCleanup)
	if err != nil {
		compensations.CompensateAndPanic(ctx, err)
	}

	services.Set(svc, services.KeyWorkflows, workflows)

	as := TestApplicationServices{
		Services: svc,
	}

	appConfig := TestApplicationConfig{
		Services:              as.Services,
		DurableExecutorConfig: config,
	}

	testApp, appCleanup := CreateApp(ctx, appConfig)
	compensations.AddCompensation(appCleanup)

	cleanup := compensations.Compensate

	return testApp, cleanup
}
