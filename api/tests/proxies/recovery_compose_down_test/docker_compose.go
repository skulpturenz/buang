package proxiesrecoverycomposedown

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type dockerComposeProxyImpl struct {
	mu        sync.Mutex
	downCount int
}

func (c *dockerComposeProxyImpl) NewComposeService(cli command.Cli, options ...compose.Option) (api.Compose, error) {
	svc, err := compose.NewComposeService(cli, options...)
	if err != nil {
		return nil, err
	}

	return &proxyDockerCompose{
		Compose: svc,
		factory: c,
	}, nil
}

var _ api.Compose = (*proxyDockerCompose)(nil)

type proxyDockerCompose struct {
	api.Compose
	factory *dockerComposeProxyImpl
}

func (d *proxyDockerCompose) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	d.factory.mu.Lock()
	defer d.factory.mu.Unlock()

	initialCount := d.factory.downCount
	d.factory.downCount++

	if initialCount < 2 {
		return fmt.Errorf("error!!!!!!")
	}

	return d.Compose.Down(ctx, projectName, options)
}
