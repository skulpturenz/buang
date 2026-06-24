package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type GetProjectWebhooksParams struct {
	ProjectID int64
}

type GetProjectWebhooksResult struct {
	Webhooks []interfaces.ProjectWebhook
}

func (p GetProjectWebhooksParams) Exec(ctx context.Context, s *app.ApplicationServices) (*GetProjectWebhooksResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.GetProjectWebhooks(ctx, interfaces.GetProjectWebhooksParams{
		ProjectID: p.ProjectID,
	})
	if err != nil {
		return nil, err
	}

	ret := GetProjectWebhooksResult{
		Webhooks: result,
	}

	return &ret, nil
}