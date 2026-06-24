package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) GetProjectWebhooks(ctx context.Context, args interfaces.GetProjectWebhooksParams) ([]interfaces.ProjectWebhook, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.GetProjectWebhooks(ctx, args.ProjectID)

	ret := []interfaces.ProjectWebhook{}
	for _, x := range res {
		ret = append(ret, ProjectWebhook(x))
	}

	return ret, err
}