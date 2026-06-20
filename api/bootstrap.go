package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"
	"strings"
)

//go:embed bootstrap/docker-compose.prod.yml
var bootstrapFiles embed.FS

func bootstrap(ctx context.Context, s *app.ApplicationServices) error { // TODO: test
	bd, err := os.MkdirTemp("/var/tmp", "buang-bootstrap-*")
	if err != nil {
		return err
	}

	data, err := bootstrapFiles.ReadFile("bootstrap/docker-compose.prod.yml")
	if err != nil {
		return fmt.Errorf("unable to find bootstrap resources")
	}

	projectName := "buang"
	configPath := filepath.Join(bd, "./docker-compose.prod.yml")
	envs := map[string]any{
		"GO_ENV":                  "production",
		"BUANG_BOOTSTRAP_PROJECT": projectName,
		"BUANG_BOOTSTRAP_DIR":     bd,
	}
	for _, v := range os.Environ() {
		env := strings.Split(v, "=")
		if len(env) == 2 {
			k := env[0]
			v := env[1]

			envs[k] = v
		}
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	upParams := docker.ComposeUpParams{
		ProjectName: projectName,
		ConfigPaths: []string{
			configPath,
		},
		Environment: envs,
	}

	res, _, err := upParams.Exec(ctx, s)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, fmt.Sprintf("buang bootstrapped! project: %v, services: %v", res.ProjectName, strings.Join(res.ServiceNames, ", ")))

	return nil
}
