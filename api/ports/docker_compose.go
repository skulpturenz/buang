package ports

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type DockerCompose interface {
	NewComposeService(cli command.Cli, options ...compose.Option) (api.Compose, error)
}

type DockerComposeImpl struct{}

func (c *DockerComposeImpl) NewComposeService(cli command.Cli, options ...compose.Option) (api.Compose, error) {
	svc, err := compose.NewComposeService(cli, options...)
	if err != nil {
		return nil, err
	}

	return svc, nil
}
