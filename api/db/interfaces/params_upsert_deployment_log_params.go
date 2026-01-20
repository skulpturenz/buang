package interfaces

type UpsertDeploymentLogParams struct {
	ProjectID    int64
	DeploymentID int64
	Log          *string
}
