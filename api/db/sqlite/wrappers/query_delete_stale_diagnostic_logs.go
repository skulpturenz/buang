package wrappers

import (
	"context"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) DeleteStaleDiagnosticLogs(ctx context.Context) error {
	m := sqlitemodels.Queries(q)
	err := m.DeleteStaleDiagnosticLogs(ctx)
	if err != nil {
		return err
	}

	return err
}
