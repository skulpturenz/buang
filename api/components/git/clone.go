package git

import (
	"context"
	"os"
	"reflect"
	"skulpture/buang/app"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/negrel/assert"
)

type CloneParams struct {
	URL               string
	RecurseSubmodules int
	Branch            string `validate:"required_without=Sha,excluded_with=Sha"`
	Hash              string `validate:"required_without=Branch,excluded_with=Branch"`
	Username          *string
	Password          *string
}

type CloneResult struct {
	URL  string
	Hash string
	Dir  string
}

func (c CloneParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CloneResult, func(ctx context.Context), error) {
	assert.True(!reflect.ValueOf(c.Hash).IsZero() && !reflect.ValueOf(c.Branch).IsZero(), "hash and branch are mutually exclusive")

	dir, err := os.MkdirTemp("", "buang-*")
	cleanup := func(ctx context.Context) {
		os.RemoveAll(dir)
	}
	assert.DirExists(dir, "invalid dir")

	if err != nil {
		return nil, cleanup, err
	}

	opts := git.CloneOptions{
		URL:               c.URL,
		NoCheckout:        true,
		InsecureSkipTLS:   true,
		RecurseSubmodules: git.SubmoduleRecursivity(c.RecurseSubmodules),
	}
	if c.Username != nil && c.Password != nil {
		opts.Auth = &http.BasicAuth{
			Username: *c.Username,
			Password: *c.Password,
		}
	}

	r, err := git.PlainCloneContext(ctx, dir, &opts)
	if err != nil {
		return nil, cleanup, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, cleanup, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(c.Branch),
		Hash:   plumbing.NewHash(c.Hash),
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
