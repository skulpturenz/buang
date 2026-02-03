package workers

import (
	"context"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
)

type TemporalConfig struct {
	Services workersshared.WorkflowServices

	HostPort *string
}

func (c TemporalConfig) New(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
	wc := workerConfig{
		Type:     enumsdurableexecutors.Temporal,
		services: c.Services,

		TemporalHostPort: c.HostPort,
	}

	return wc.new(ctx)
}
