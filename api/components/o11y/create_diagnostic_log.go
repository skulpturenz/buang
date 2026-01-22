package o11y

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
)

type CreateDiagnosticLogParams struct {
	Type enumsdiagnosticlogtype.DiagnosticLogType
	Log  map[string]any
}

type CreateDiagnosticLogResult struct {
	ID int64
}

func (p *CreateDiagnosticLogParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CreateDiagnosticLogResult, error) {
	q := *s.Queries

	res, err := q.CreateDiagnosticLog(ctx, interfaces.CreateDiagnosticLogParams{
		Type: p.Type,
		Log:  p.Log,
	})
	if err != nil {
		return nil, err
	}

	ret := CreateDiagnosticLogResult{
		ID: res.GetId(),
	}

	return &ret, nil
}
