package docker

import (
	"context"
	"io"
	"skulpture/buang/app"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type StopContainerParams struct {
	ID     string
	Writer io.Writer
}

type StopContainerResult struct {
	ID string
}

func (c StopContainerParams) StopContainer(ctx context.Context, s *app.ApplicationServices) (*StopContainerResult, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	if err := cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
	}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, c.ID, container.WaitConditionRemoved)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, c.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, err
	}
	if c.Writer != nil {
		io.Copy(c.Writer, out)
	}
	if c.Writer == nil {
		io.Copy(io.Discard, out)
	}

	return &StopContainerResult{
		ID: c.ID,
	}, nil
}
