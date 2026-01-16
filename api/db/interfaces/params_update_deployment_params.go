package interfaces

import "time"

type UpdateDeploymentParams struct {
	ID         int64
	Url        *string
	Status     int16
	DeployedAt *time.Time
}
