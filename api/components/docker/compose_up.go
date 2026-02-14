package docker

import (
	"context"
	"fmt"
	"io"
	"os"
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
	os.Setenv("DOCKER_BUILDKIT", "1")
	cliOptions := []command.CLIOption{
		command.WithAPIClient(s.Docker),
	}
	if c.Writer != nil {
		cliOptions = append(cliOptions, command.WithCombinedStreams(c.Writer))
	}

	cli, err := command.NewDockerCli(cliOptions...)
	if err != nil {
		return nil, nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, nil, err
	}

	options := []compose.Option{
		compose.WithPrompt(compose.AlwaysOkPrompt()),
		compose.WithContextInfo(&contextInfo{cli: cli}),
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

	svc, err := s.Ports.DockerCompose().NewComposeService(cli, options...)
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

	projectWithServiceAliases, err := project.WithServicesTransform(func(name string, s types.ServiceConfig) (types.ServiceConfig, error) {
		alias := fmt.Sprintf("%v_%v", project.Name, name)
		updatedService := s

		if len(s.Networks) == 0 {
			updatedService.Networks = map[string]*types.ServiceNetworkConfig{
				"default": {
					Aliases: []string{alias},
				},
			}
		} else {
			for networkName, networkConfig := range updatedService.Networks {
				updatedNetworkConfig := networkConfig
				if updatedNetworkConfig == nil {
					updatedNetworkConfig = &types.ServiceNetworkConfig{}
				}

				updatedNetworkConfig.Aliases = append(updatedNetworkConfig.Aliases, alias)
				updatedService.Networks[networkName] = updatedNetworkConfig
			}
		}

		return updatedService, nil
	})

	logConsumer := logConsumer{
		Writer: c.Writer,
	}
	logCtx, cancelLogCtx := context.WithCancel(ctx)
	go followSvcLogs(logCtx, projectWithServiceAliases.Name, logConsumer, svc)
	defer cancelLogCtx()

	buildOptions := api.BuildOptions{
		Pull: true,
		Push: true,
		Deps: true,
	}
	if c.Writer != nil {
		buildOptions.Out = c.Writer
	}

	err = svc.Up(ctx, projectWithServiceAliases, api.UpOptions{
		Create: api.CreateOptions{
			Build:         &buildOptions,
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
			ProjectName: projectWithServiceAliases.Name,
			ConfigPaths: c.ConfigPaths,
		}

		down.Exec(ctx, s)
	}

	return &ComposeUpResult{
		ProjectName:  projectWithServiceAliases.Name,
		ConfigPaths:  c.ConfigPaths,
		Environment:  projectWithServiceAliases.Environment.Values(),
		ServiceNames: projectWithServiceAliases.ServiceNames(),
	}, cleanup, nil
}

func followSvcLogs(ctx context.Context, projectName string, logConsumer logConsumer, svc api.Compose) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := svc.Logs(ctx, projectName, logConsumer, api.LogOptions{ // blocks indefinitely
				Follow:     true,
				Timestamps: true,
			})

			if err != nil {
				return
			}
		}
	}
}

type contextInfo struct {
	cli command.Cli
}

func (c *contextInfo) CurrentContext() string {
	return c.cli.CurrentContext()
}

func (c *contextInfo) ServerOSType() string {
	return c.cli.ServerInfo().OSType
}

func (c *contextInfo) BuildKitEnabled() (bool, error) {
	// the cli checks env by default
	// override and hardcode to true
	// see: dockerCliContextInfo.BuildKitEnabled
	// DockerCli.BuildKitEnabled
	return true, nil
}
