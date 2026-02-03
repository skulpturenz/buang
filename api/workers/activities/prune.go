package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"
)

type Prune app.ApplicationServices

type PruneParams struct{}

type PruneResult struct{}

func (p *Prune) Prune(ctx context.Context) (*PruneResult, error) {
	s := app.ApplicationServices(*p)

	params := docker.PruneParams{}

	_, err := params.Prune(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &PruneResult{}, nil
}
