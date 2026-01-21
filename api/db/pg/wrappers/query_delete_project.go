package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) DeleteProject(ctx context.Context, projectID int64) (interfaces.Project, error) {
	m := pgmodels.Queries(q)
	ret, err := m.DeleteProject(ctx, projectID)

	return Project(ret), err
}
