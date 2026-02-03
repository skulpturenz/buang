package wrappers

import (
	"context"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) DeleteStaleDiagnosticLogs(ctx context.Context) error {
	m := pgmodels.Queries(q)
	err := m.DeleteStaleDiagnosticLogs(ctx)
	if err != nil {
		return err
	}

	return err
}
