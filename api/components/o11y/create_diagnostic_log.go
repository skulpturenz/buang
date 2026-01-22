package o11y

import (
	"context"
	"log/slog"
	"skulpture/buang/db/interfaces"
)

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
