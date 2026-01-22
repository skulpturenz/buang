package o11y

import (
	"context"
	"fmt"
	"log/slog"
)

func PruneStaleDiagnosticLogs(ctx context.Context) error {
	i := getInstance()
	if i == nil {
		err := fmt.Errorf("diagnostic log singleton not initialized, nothing will be captured")
		slog.Error(err.Error())

		return err
	}

	q := *i.q

	err := q.DeleteStaleDiagnosticLogs(ctx)
	if err != nil {
		return err
	}

	return nil
}
