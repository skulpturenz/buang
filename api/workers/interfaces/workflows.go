package workersinterfaces

import "context"

type Workflows interface {
	BuangBranch(ctx context.Context, params BuangBranchParams) error
	BuangDeployment(ctx context.Context, params BuangDeploymentParams) error
	CreateDeployment(ctx context.Context, params CreateDeploymentParams) error
	HasDeployed(ctx context.Context, params HasDeployedParams) error
}
