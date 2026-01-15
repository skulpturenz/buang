package interfaces

import "context"

type Queries interface {
	SelectAllProjectsDesc(ctx context.Context) ([]Project, error)
}
