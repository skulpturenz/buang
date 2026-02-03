package wrappers

import (
	"context"
	"encoding/json"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) CreateDiagnosticLog(ctx context.Context, arg interfaces.CreateDiagnosticLogParams) (interfaces.DiagnosticLog, error) {
	m := sqlitemodels.Queries(q)

	log, err := json.Marshal(arg.Log)
	if err != nil {
		return nil, err
	}

	ret, err := m.CreateDiagnosticLog(ctx, sqlitemodels.CreateDiagnosticLogParams{
		Type: arg.Type.String(),
		Log:  log,
	})

	return DiagnosticLog(ret), err
}
