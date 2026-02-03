package ports

import (
	"context"

	"github.com/go-git/go-git/v6"
)

type Git interface {
	PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error)
}

type GitImpl struct{}

func (c *GitImpl) PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainCloneContext(ctx, path, o)
}
