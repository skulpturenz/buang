package workersdbosclient

import (
	"context"

	"skulpture/buang/app"
	workersdbos "skulpture/buang/workers/dbos"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type DbosWorkflow[P any, R any] func(s app.ApplicationServices, c dbos.DBOSContext)

type DbosConfig struct {
	Services workersshared.WorkflowServices

	AppName         string
	DatabaseURL     string
	ConductorAPIKey *string
	ConductorURL    *string
	AdminServerPort *int
}

type Client struct {
	ctx dbos.DBOSContext
	s   app.ApplicationServices
}

func New(ctx context.Context, cfg DbosConfig) (workersinterfaces.Workflows, error) {
	c := dbos.Config{
		AppName:     cfg.AppName,
		DatabaseURL: cfg.DatabaseURL,
	}
	if cfg.ConductorAPIKey != nil {
		c.ConductorAPIKey = *cfg.ConductorAPIKey
	}
	if cfg.ConductorURL != nil {
		c.ConductorURL = *cfg.ConductorURL
	}
	if cfg.AdminServerPort != nil {
		c.AdminServer = true
		c.AdminServerPort = *cfg.AdminServerPort
	}

	dbosContext, err := dbos.NewDBOSContext(ctx, c)
	if err != nil {
		return nil, err
	}

	client := &Client{ctx: dbosContext, s: cfg.Services.ToAppServices()}

	workflows := []DbosWorkflow[any, any]{
		workersdbos.Deployment,
		workersdbos.Buang,
		workersdbos.Housekeping,
	}

	for _, w := range workflows {
		w(cfg.Services.ToAppServices(), dbosContext)
	}

	err = dbos.Launch(dbosContext)
	if err != nil {
		return nil, err
	}

	return client, nil
}
