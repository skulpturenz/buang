package activities

import (
	"context"
	"os"
	"path/filepath"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"
	"strings"
)

type AutoupdateBuang app.ApplicationServices

type AutoupdateBuangResult struct{}

func (bd *AutoupdateBuang) AutoupdateBuang(ctx context.Context) (*AutoupdateBuangResult, error) {
	s := app.ApplicationServices(*bd)

	projectName, ok := os.LookupEnv("BUANG_BOOTSTRAP_PROJECT")
	if !ok {
		return nil, nil
	}

	bootstrapDir, ok := os.LookupEnv("BUANG_BOOTSTRAP_DIR")
	if !ok {
		return nil, nil
	}

	envs := map[string]any{
		"GO_ENV": "production",
	}
	for _, v := range os.Environ() {
		env := strings.Split(v, "=")
		if len(env) == 2 {
			k := env[0]
			v := env[1]

			envs[k] = v
		}
	}

	upParams := docker.ComposeUpParams{
		ProjectName: projectName,
		ConfigPaths: []string{filepath.Join(bootstrapDir, "./docker-compose.prod.yml")},
		Environment: envs,
	}

	_, _, err := upParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &AutoupdateBuangResult{}, nil
}
