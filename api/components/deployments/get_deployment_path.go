package deployments

import (
	"fmt"
	"path/filepath"
	"time"
)

const TRAEFIK_DYNAMIC_CONFIG = "/app/deployments"

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
