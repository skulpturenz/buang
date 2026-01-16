package activities

import (
	"context"
	"os"
	"skulpture/buang/app"
)

type CreateDynamicConfigDir app.ApplicationServices

type CreateDynamicConfigDirParams struct{}

type CreateDynamicConfigDirResult struct{}

func (cdcd *CreateDynamicConfigDir) Exec(ctx context.Context, c CreateDynamicConfigDirParams) error {
	err := os.MkdirAll(TRAEFIK_DYNAMIC_CONFIG, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
