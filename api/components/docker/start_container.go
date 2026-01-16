package docker

import (
	"context"
	"io"
	"skulpture/buang/app"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type StartContainerParams struct {
	Name         string
	Image        string
	Writer       io.Writer
	Labels       map[string]string
	Env          []string
	Cmd          []string
	ExposedPorts nat.PortSet
	PortBindings nat.PortMap
}

type StartContainerResult struct {
	ID string
}

func (c StartContainerParams) StartContainer(ctx context.Context, s *app.ApplicationServices) (*StartContainerResult, func(ctx context.Context), error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, c.Image, image.PullOptions{})
	if err != nil {
		return nil, nil, err
	}
	if c.Writer != nil {
		io.Copy(c.Writer, reader)
	}
	if c.Writer == nil {
		io.Copy(io.Discard, reader)
	}

	containerConfig := &container.Config{
		Image:        c.Image,
		Cmd:          []string{},
		Labels:       c.Labels,
		Env:          c.Env,
		ExposedPorts: c.ExposedPorts,
	}
	hostConfig := &container.HostConfig{
		PortBindings: c.PortBindings,
	}

	create, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		c.Name,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := cli.ContainerStart(ctx, create.ID, container.StartOptions{}); err != nil {
		return nil, nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, create.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, nil, err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, create.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, nil, err
	}
	if c.Writer != nil {
		io.Copy(c.Writer, out)
	}
	if c.Writer == nil {
		io.Copy(io.Discard, out)
	}

	cleanup := func(ctx context.Context) {
		down := StopContainerParams{
			ID:     create.ID,
			Writer: c.Writer,
		}

		down.StopContainer(ctx, s)
	}

	return &StartContainerResult{
		ID: create.ID,
	}, cleanup, nil
}
