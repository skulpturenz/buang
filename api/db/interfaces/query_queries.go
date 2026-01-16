package interfaces

import "context"

type Queries interface {
	SelectProjectsDesc(ctx context.Context, arg SelectProjectsDescParams) ([]Project, error)
	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
}
