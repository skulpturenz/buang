package docker

import (
	"context"
	"io"
	"skulpture/buang/app"

	"github.com/docker/docker/api/types/container"
)

type StopContainerParams struct {
	ID     string
	Writer io.Writer
}

type StopContainerResult struct {
	ID string
}

func (c StopContainerParams) StopContainer(ctx context.Context, s *app.ApplicationServices) (*StopContainerResult, error) {
	if err := s.Docker.ContainerRemove(ctx, c.ID, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
	}); err != nil {
		return nil, err
	}

	statusCh, errCh := s.Docker.ContainerWait(ctx, c.ID, container.WaitConditionRemoved)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := s.Docker.ContainerLogs(ctx, c.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
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
