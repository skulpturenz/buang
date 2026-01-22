package o11y

import (
	"context"
	"log/slog"
	"skulpture/buang/db/interfaces"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"sync"
)

type diagnosticLog struct {
	q *interfaces.Queries
}

type CreateDiagnosticLogParams struct {
	Type enumsdiagnosticlogtype.DiagnosticLogType
	Log  map[string]any
}

type CreateDiagnosticLogResult struct {
	ID int64
}

var instance *diagnosticLog

var once sync.Once

func NewS(q *interfaces.Queries) *diagnosticLog {
	once.Do(func() {
		instance = &diagnosticLog{
			q: q,
		}
	})

	return instance
}

// ignoring errors here technically breaks the rules of workflows
// but its so that we can capture the information we need
func getInstance() *diagnosticLog {
	if instance == nil {
		slog.Error("diagnostic log singleton not initialized, nothing will be captured")

		return nil
	}

	return instance
}

// ignoring errors here technically breaks the rules of workflows
// but its so that we can capture the information we need
func (p *CreateDiagnosticLogParams) Exec(ctx context.Context) *CreateDiagnosticLogResult {
	i := getInstance()
	if i == nil {
		slog.Error("diagnostic log singleton not initialized, nothing will be captured")

		return &CreateDiagnosticLogResult{}
	}

	q := *i.q

	res, _ := q.CreateDiagnosticLog(ctx, interfaces.CreateDiagnosticLogParams{
		Type: p.Type,
		Log:  p.Log,
	})

	ret := CreateDiagnosticLogResult{
		ID: res.GetId(),
	}

	return &ret
}
