package workflowwrappers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type HasDeployedParams struct {
	ProjectId    int64
	DeploymentId int64
	Block        bool
}

func (p HasDeployedParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("has deployed panic: %v", r))

			log := map[string]any{
				"stack": string(debug.Stack()),
			}

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  log,
			}
			p.Exec(ctx)

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("has deployed panic")
			}
		}
	}()

	executor, durableExecutor := s.GetDurableExecutor()

	cdp := CreateDeploymentParams{
		ProjectId:    p.ProjectId,
		DeploymentId: p.DeploymentId,
	}
	workflowId := cdp.GetWorkflowId()
	if durableExecutor == enumsdurableexecutors.Temporal {
		run := executor.(app.TemporalClient).GetWorkflow(ctx, workflowId, "")
		if err != nil {
			return err
		}

		err = run.Get(ctx, nil)
		if err != nil {
			return err
		}

		return nil
	} else {
		handle, err := dbos.RetrieveWorkflow[bool](executor.(app.DbosContext), cdp.GetWorkflowId())
		if err != nil {
			return err
		}

		_, err = handle.GetResult()
		if err != nil {
			return err
		}

		return nil
	}
}
