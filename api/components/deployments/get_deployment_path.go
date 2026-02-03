package deployments

import (
	"fmt"
	"path/filepath"
	constantsenvs "skulpture/buang/constants/envs"
	"time"
)

var TRAEFIK_DYNAMIC_CONFIG = constantsenvs.BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR.Value()

type GetDeploymentPathParams struct {
	ProjectId    int64
	DeploymentId int64
	Branch       string
	Sha          string
	DeployedAt   time.Time
}

func GetDeploymentPath(p GetDeploymentPathParams) string {
	projectName := GetProjectName(GetProjectNameParams{
		ProjectId:    p.ProjectId,
		DeploymentId: p.DeploymentId,
		Branch:       p.Branch,
		Sha:          p.Sha,
		DeployedAt:   p.DeployedAt,
	})

	return filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("buang-%v.yaml", projectName))
}
