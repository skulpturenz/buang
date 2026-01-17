package docker

import (
	"context"
	"fmt"
	"io"
	"skulpture/buang/app"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type ComposeUpParams struct {
	ConfigPaths    []string
	ProjectName    string
	Writer         io.Writer
	DryRun         bool
	EventProcessor api.EventProcessor
	Environment    map[string]any
}

type ComposeUpResult struct {
	ProjectName  string
	ConfigPaths  []string
	Environment  []string
	ServiceNames []string
}

func (c ComposeUpParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ComposeUpResult, func(ctx context.Context), error) {
	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, nil, err
	}

	options := []compose.Option{}
	if c.Writer != nil {
		options = append(options, compose.WithOutputStream(c.Writer), compose.WithErrorStream(c.Writer))
	}
	if c.DryRun == true {
		options = append(options, compose.WithDryRun)
	}
	if c.EventProcessor != nil {
		options = append(options, compose.WithEventProcessor(c.EventProcessor))
	}

	svc, err := compose.NewComposeService(cli, options...)
	if err != nil {
		return nil, nil, err
	}

	project, err := svc.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: c.ConfigPaths,
		ProjectName: c.ProjectName,
	})
	if err != nil {
		return nil, nil, err
	}
	envVars := []string{}
	for k, v := range c.Environment {
		envVars = append(envVars, fmt.Sprintf("%v=%v", k, v))
	}

	project.Environment.Merge(types.NewMapping(envVars))

	err = svc.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{
			Build: &api.BuildOptions{
				Pull: true,
				Deps: true,
			},
			RemoveOrphans: true,
		},
		Start: api.StartOptions{
			Wait:        true,
			WaitTimeout: time.Minute * 10,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func(ctx context.Context) {
		down := ComposeDownParams{
			ProjectName: project.Name,
			ConfigPaths: c.ConfigPaths,
		}

		down.Exec(ctx, s)
	}

	return &ComposeUpResult{
		ProjectName:  project.Name,
		ConfigPaths:  c.ConfigPaths,
		Environment:  project.Environment.Values(),
		ServiceNames: project.ServiceNames(),
	}, cleanup, nil
}
