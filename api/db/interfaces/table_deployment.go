package interfaces

type Deployment interface {
	GetId() int64
	GetRepositoryId() int64
	GetUrl() *string
	GetStatus() int16
	GetSha() *string
}
