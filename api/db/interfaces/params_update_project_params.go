package interfaces

type UpdateProjectParams struct {
	Repository    string
	RequiresAuthn bool
	Username      *string
	Password      *string
	ComposePath   string
	ProjectID     int64
}
