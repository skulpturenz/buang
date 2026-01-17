package interfaces

type CreateDeploymentParams struct {
	ProjectID         int64
	Sha               *string
	Status            int16
	ServiceEntrypoint string
}
