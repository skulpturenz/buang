package o11y

import (
	"context"
	"log/slog"
	"skulpture/buang/db/interfaces"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
)

type CreateDiagnosticLogParams struct {
	Type enumsdiagnosticlogtype.DiagnosticLogType
	Log  any
}

type CreateDiagnosticLogResult struct {
	ID int64
}

func (p *CreateDiagnosticLogParams) Exec(ctx context.Context) *CreateDiagnosticLogResult {
	i := getInstance()
	if i == nil {
		slog.Error("diagnostic log singleton not initialized, nothing will be captured")

		return nil
	}

	q := *i.q

	log := map[string]any{}
	switch p.Type {
	case enumsdiagnosticlogtype.Panic:
		log["panic"] = p.Log
	case enumsdiagnosticlogtype.DockerStats:
		log["docker_stats"] = p.Log
	default:
		return nil
	}

	res, _ := q.CreateDiagnosticLog(ctx, interfaces.CreateDiagnosticLogParams{
		Type: p.Type,
		Log:  log,
	})

	ret := CreateDiagnosticLogResult{
		ID: res.GetId(),
	}

	return &ret
}
