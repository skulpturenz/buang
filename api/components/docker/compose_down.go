package docker

import (
	"context"
	"io"
	"skulpture/buang/app"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/negrel/assert"
)

type ComposeDownParams struct {
	ProjectName string
	ConfigPaths []string
	Writer      io.Writer
}

type ComposeDownResult struct{}

func (c ComposeDownParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ComposeDownResult, error) {
	docker, ok := s.GetDocker()
	assert.True(ok, "docker service not found")
	ports, ok := s.GetPorts()
	assert.True(ok, "ports service not found")

	cliOptions := []command.CLIOption{
		command.WithAPIClient(docker),
	}
	if c.Writer != nil {
		cliOptions = append(cliOptions, command.WithCombinedStreams(c.Writer))
	}

	cli, err := command.NewDockerCli(cliOptions...)
	if err != nil {
		return nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, err
	}

	svc, err := ports.DockerCompose().NewComposeService(cli)
	if err != nil {
		return nil, err
	}

	project, err := svc.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: c.ConfigPaths,
		ProjectName: c.ProjectName,
	})
	if err != nil {
		return nil, err
	}

	timeout := time.Minute * 2
	err = svc.Down(ctx, project.Name, api.DownOptions{
		RemoveOrphans: true,
		Timeout:       &timeout,
		Volumes:       true,
		Images:        "all",
	})
	if err != nil {
		return nil, err
	}

	return &ComposeDownResult{}, nil
}
