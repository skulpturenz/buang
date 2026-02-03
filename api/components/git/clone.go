package git

import (
	"context"
	"fmt"
	"io"
	"os"
	"skulpture/buang/app"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/negrel/assert"
)

type CloneParams struct {
	URL               string
	RecurseSubmodules int
	Branch            string
	Hash              string
	Username          *string
	Password          *string
	Writer            io.Writer
}

type CloneResult struct {
	URL  string
	Hash string
	Dir  string
}

func (c *CloneParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CloneResult, func(ctx context.Context), error) {
	dir, err := os.MkdirTemp("/var/tmp", "buang-*")
	cleanup := func(ctx context.Context) {
		os.RemoveAll(dir)
	}
	assert.DirExists(dir, "invalid dir")

	if err != nil {
		return nil, cleanup, err
	}

	opts := git.CloneOptions{
		URL:               c.URL,
		NoCheckout:        false,
		InsecureSkipTLS:   true,
		ReferenceName:     plumbing.NewHashReference(plumbing.ReferenceName(c.Branch), plumbing.NewHash(c.Hash)).Name(),
		RecurseSubmodules: git.SubmoduleRecursivity(c.RecurseSubmodules),
		Progress:          c.Writer,
		SingleBranch:      true,
	}
	if c.Username != nil && c.Password != nil {
		opts.Auth = &http.BasicAuth{
			Username: *c.Username,
			Password: *c.Password,
		}
	}

	r, err := s.Ports.Git().PlainCloneContext(ctx, dir, &opts)
	if err != nil {
		return nil, cleanup, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, cleanup, err
	}

	hash, ok := plumbing.FromHex(c.Hash)
	if !ok {
		return nil, cleanup, fmt.Errorf("unable to clone project from hash %v", c.Hash)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})
	if err != nil {
		return nil, cleanup, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, cleanup, err
	}

	_, err = fmt.Fprintf(c.Writer, "Checked out commit %v, branch %v", head.Hash(), c.Branch)
	if err != nil {
		return nil, cleanup, err
	}

	return &CloneResult{
		URL:  c.URL,
		Hash: c.Hash,
		Dir:  dir,
	}, cleanup, nil
}
