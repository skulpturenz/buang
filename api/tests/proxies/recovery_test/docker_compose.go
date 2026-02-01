package proxiesrecoverytest

import (
	"context"
	"fmt"
	"sync"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type dockerComposeProxyImpl struct {
	mu      sync.Mutex
	upCount int
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

func (d *proxyDockerCompose) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	d.factory.mu.Lock()
	defer d.factory.mu.Unlock()

	initialCount := d.factory.upCount
	d.factory.upCount++

	if initialCount < 2 {
		return fmt.Errorf("ahhhH!!! error!!!! too many image pull requests!!!")
	}

	return d.Compose.Up(ctx, project, options)
}

func (d *proxyDockerCompose) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	return d.Compose.Down(ctx, projectName, options)
}
