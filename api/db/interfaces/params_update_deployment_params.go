package interfaces

import "time"

type UpdateDeploymentParams struct {
	Url        *string
	Status     int16
	DeployedAt *time.Time
	ClonePath  *string
	ID         int64
	ProjectID  int64
}
