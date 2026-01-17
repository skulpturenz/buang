package interfaces

type CreateProjectParams struct {
	Repository    string
	RequiresAuthn bool
	Username      *string
	Password      *string
	ComposePath   string
}
