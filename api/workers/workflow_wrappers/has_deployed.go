package workflowwrappers

import (
	"context"
	"skulpture/buang/app"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type HasDeployedParams struct {
	ProjectId    int64
	DeploymentId int64
	Block        bool
}

func (p HasDeployedParams) Exec(ctx context.Context, s app.ApplicationServices) (err error) {
	defer RecoverWorkflowPanic("HasDeployed", err)

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
