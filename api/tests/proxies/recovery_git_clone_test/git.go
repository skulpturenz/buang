package proxiesrecoverygitclone

import (
	"context"
	"fmt"
	"sync"

	"skulpture/buang/ports"

	"github.com/go-git/go-git/v6"
)

type gitProxyImpl struct {
	mu         sync.Mutex
	cloneCount int
}

func (c *gitProxyImpl) PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	initialCount := c.cloneCount
	c.cloneCount++

	if initialCount < 2 {
		return nil, fmt.Errorf("git clone failed")
	}

	return (&ports.GitImpl{}).PlainCloneContext(ctx, path, o)
}
