package interfaces

type DeploymentLog interface {
	GetId() int64
	GetDeploymentId() int64
	GetLog() *string
}
