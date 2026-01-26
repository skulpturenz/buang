package workersinterfaces

import "fmt"

type CreateDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
	Block        bool
}

func (p CreateDeploymentParams) GetWorkflowId() string {
	return fmt.Sprintf("create-project-%v-deployment-%v", p.ProjectId, p.DeploymentId)
}
