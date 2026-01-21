package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) DeleteProject(ctx context.Context, projectID int64) (interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.DeleteProject(ctx, projectID)

	return Project(ret), err
}
