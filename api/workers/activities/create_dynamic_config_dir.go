package activities

import (
	"context"
	"os"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
)

type CreateDynamicConfigDir app.ApplicationServices

type CreateDynamicConfigDirParams struct{}

type CreateDynamicConfigDirResult struct{}

func (cdcd *CreateDynamicConfigDir) CreateDynamicConfigDir(ctx context.Context, c CreateDynamicConfigDirParams) error {
	err := os.MkdirAll(deployments.TRAEFIK_DYNAMIC_CONFIG, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
