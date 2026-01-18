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
	"github.com/docker/docker/client"
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	imagesPruneReport, err := cli.ImagesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	containersPruneReport, err := cli.ContainersPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	volumesPruneReport, err := cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, err
	}

	buildCachePruneReport, err := cli.BuildCachePrune(ctx, build.CachePruneOptions{All: true})
	if err != nil {
		return nil, err
	}

	networksPruneReport, err := cli.NetworksPrune(ctx, filters.NewArgs())
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
