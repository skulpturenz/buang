package git

import (
	"context"
	"os"
	"skulpture/buang/app"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type CloneParams struct {
	URL               string
	Depth             int
	RecurseSubmodules int
	Hash              string
}

type CloneResult struct {
	URL  string
	Hash string
	Dir  string
}

func (c CloneParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CloneResult, func(ctx context.Context), error) {
	dir, err := os.MkdirTemp("", "buang-*")
	cleanup := func(ctx context.Context) {
		os.RemoveAll(dir)
	}

	if err != nil {
		return nil, cleanup, err
	}

	r, err := git.PlainCloneContext(ctx, dir, &git.CloneOptions{
		URL:               c.URL,
		NoCheckout:        true,
		InsecureSkipTLS:   true,
		Depth:             c.Depth,
		RecurseSubmodules: git.SubmoduleRecursivity(c.RecurseSubmodules),
	})
	if err != nil {
		return nil, cleanup, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, cleanup, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(c.Hash),
	})
	if err != nil {
		return nil, cleanup, err
	}

	return &CloneResult{
		URL:  c.URL,
		Hash: c.Hash,
		Dir:  dir,
	}, cleanup, nil
}
