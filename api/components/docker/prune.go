package docker

import (
	"context"
	"skulpture/buang/app"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/negrel/assert"
)

type PruneParams struct {
}

type PruneResult struct {
	ImagesPruneReport     image.PruneReport
	ContainersPruneReport container.PruneReport
	VolumesPruneReport    volume.PruneReport
	BuildCachePruneReport *build.CachePruneReport
	NetworkPruneReport    network.PruneReport
}

func (c PruneParams) Prune(ctx context.Context, s *app.ApplicationServices) (*PruneResult, error) {
	docker, ok := s.GetDocker()
	assert.True(ok, "docker service not found")

	imagesPruneReport, err := docker.ImagesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	containersPruneReport, err := docker.ContainersPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	volumesPruneReport, err := docker.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	buildCachePruneReport, err := docker.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		return nil, err
	}

	networksPruneReport, err := docker.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	return &PruneResult{
		ImagesPruneReport:     imagesPruneReport,
		ContainersPruneReport: containersPruneReport,
		VolumesPruneReport:    volumesPruneReport,
		BuildCachePruneReport: buildCachePruneReport,
		NetworkPruneReport:    networksPruneReport,
	}, nil
}
