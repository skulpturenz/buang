package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
)

type BuangDeploymentParams struct {
	ID int64
}

type BuangDeploymentResult struct{}

func (d BuangDeploymentParams) Exec(ctx context.Context, s *app.ApplicationServices) (*BuangDeploymentResult, error) {
	q := *s.Queries

	_, err := q.UpdateDeploymentStatus(ctx, interfaces.UpdateDeploymentStatusParams{
		ID:     d.ID,
		Status: int16(enumsdeploymentstatus.Buang),
	})
	if err != nil {
		return nil, err
	}

	ret := BuangDeploymentResult{}

	return &ret, nil
}
