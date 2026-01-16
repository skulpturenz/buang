package interfaces

import "time"

type Deployment interface {
	GetId() int64
	GetProjectId() int64
	GetUrl() *string
	GetStatus() int16
	GetSha() *string
	GetDeployedAt() *time.Time
	GetClonePath() *string
}
