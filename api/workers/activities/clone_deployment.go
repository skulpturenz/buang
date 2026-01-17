package activities

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"skulpture/buang/components/git"
	"skulpture/buang/components/projects"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
)

type CloneDeployment app.ApplicationServices

type CloneDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

type CloneDeploymentResult struct {
	Dir string
}

func (cd *CloneDeployment) CloneDeployment(ctx context.Context, d CloneDeploymentParams) (*CloneDeploymentResult, error) {
	s := app.ApplicationServices(*cd)

	findProjectParams := projects.FindProjectByIdParams{
		ID: d.ProjectId,
	}

	p, err := findProjectParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	findDeploymentParams := deployments.FindDeploymentByIdParams{
		ID:        d.DeploymentId,
		ProjectId: p.Project.GetId(),
	}

	dply, err := findDeploymentParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	if dply.Deployment.GetStatus() != int16(enumsdeploymentstatus.New) {
		return nil, fmt.Errorf("deployment %v for project %v is invalid", d.DeploymentId, d.ProjectId)
	}

	cloneParams := git.CloneParams{
		URL:               p.Project.GetRepository(),
		Depth:             1,
		RecurseSubmodules: 3,
		Hash:              dply.Deployment.GetSha(),
	}
	if p.Project.GetRequiresAuthn() {
		cloneParams.Username = p.Project.GetUsername()
		cloneParams.Password = p.Project.GetPassword()
	}

	c, _, err := cloneParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	deployingParams := deployments.UpdateDeploymentParams{
		ID:        dply.Deployment.GetId(),
		Status:    int16(enumsdeploymentstatus.Deploying),
		ClonePath: &c.Dir,
	}

	_, err = deployingParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &CloneDeploymentResult{Dir: c.Dir}, nil
}
