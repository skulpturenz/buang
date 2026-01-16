package docker

import (
	"context"
	"skulpture/buang/app"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type ComposeDownParams struct {
	ProjectName string
	ConfigPaths []string
}

type ComposeDownResult struct{}

func (c ComposeDownParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ComposeDownResult, error) {
	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, err
	}

	svc, err := compose.NewComposeService(cli)
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
