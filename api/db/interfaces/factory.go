package interfaces

import (
	"context"
)

type DbFactory interface {
	New(ctx context.Context) (Queries, func(ctx context.Context), error)
}
