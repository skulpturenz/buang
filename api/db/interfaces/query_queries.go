package interfaces

import "context"

type Queries interface {
	SelectProjectsDesc(ctx context.Context, arg SelectProjectsDescParams) ([]Project, error)
	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	CreateDeployment(ctx context.Context, arg CreateDeploymentParams) (Deployment, error)
	UpdateDeployment(ctx context.Context, arg UpdateDeploymentParams) (Deployment, error)
	SelectDeploymentsDesc(ctx context.Context, arg SelectDeploymentsDescParams) ([]Deployment, error)
	UpdateDeploymentStatus(ctx context.Context, arg UpdateDeploymentStatusParams) (Deployment, error)
	SelectProject(ctx context.Context, id int64) (Project, error)
	SelectDeployment(ctx context.Context, arg SelectDeploymentParams) (Deployment, error)
	SelectActiveDeploymentsByBranch(ctx context.Context, arg SelectActiveDeploymentsByBranchParams) ([]Deployment, error)
	SelectStaleDeployments(ctx context.Context) ([]Deployment, error)
}
