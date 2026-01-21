package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"skulpture/buang/app"
	"time"

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

type logConsumer struct {
	Writer io.Writer
}

func (c logConsumer) Log(containerName, message string) {
	c.Writer.Write(fmt.Appendf(nil, "[%v] %v", containerName, message))
}

func (c logConsumer) Err(containerName, message string) {
	c.Writer.Write(fmt.Appendf(nil, "[%v] %v", containerName, message))
}

func (c logConsumer) Status(containerName, message string) {
	c.Writer.Write(fmt.Appendf(nil, "[%v] %v", containerName, message))
}

func (c ComposeUpParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ComposeUpResult, func(ctx context.Context), error) {
	cliOptions := []command.CLIOption{}
	if c.Writer != nil {
		cliOptions = append(cliOptions, command.WithCombinedStreams(c.Writer))
	}

	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, nil, err
	}

	options := []compose.Option{
		compose.WithPrompt(compose.AlwaysOkPrompt()),
	}
	if c.Writer != nil {
		// TODO: idk what this output stream is supposed to be but its not logs when the service is deploying
		// docker cli? might make more sense
		// but not sure how `docker-compose` streams service logs
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

	envfile, err := os.CreateTemp("/var/tmp", fmt.Sprintf("%v-env-*", c.ProjectName))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(envfile.Name())
	defer envfile.Close()

	for k, v := range c.Environment {
		if _, err := fmt.Fprintf(envfile, "%v=%v\n", k, v); err != nil {
			return nil, nil, err
		}
	}

	if err := envfile.Sync(); err != nil {
		return nil, nil, err
	}

	project, err := svc.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: c.ConfigPaths,
		ProjectName: c.ProjectName,
		EnvFiles:    []string{envfile.Name()},
	})
	if err != nil {
		return nil, nil, err
	}

	// this writes logs all at once right now
	// ideally we can pass an `io.Writer` and it just writes each time there is a new log message
	// docker pushes instead of we pull if that makes more sense
	logConsumer := logConsumer{
		Writer: c.Writer,
	}

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
		logsErr := svc.Logs(ctx, project.Name, logConsumer, api.LogOptions{Timestamps: true})
		if logsErr != nil {
			return nil, nil, errors.Join(logsErr, err)
		}

		return nil, nil, err
	}

	cleanup := func(ctx context.Context) {
		down := ComposeDownParams{
			ProjectName: project.Name,
			ConfigPaths: c.ConfigPaths,
		}

		down.Exec(ctx, s)
	}

	err = svc.Logs(ctx, project.Name, logConsumer, api.LogOptions{Timestamps: true})
	if err != nil {
		return nil, nil, errors.Join(err, err)
	}

	return &ComposeUpResult{
		ProjectName:  project.Name,
		ConfigPaths:  c.ConfigPaths,
		Environment:  project.Environment.Values(),
		ServiceNames: project.ServiceNames(),
	}, cleanup, nil
}
