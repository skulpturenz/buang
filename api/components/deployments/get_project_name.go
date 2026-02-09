package deployments

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type GetProjectNameParams struct {
	ProjectId    int64
	DeploymentId int64
	Branch       string
	Sha          string
	DeployedAt   time.Time
}

func GetProjectName(p GetProjectNameParams) string {
	sha := fmt.Sprintf("%.*s", 8, p.Sha)
	projectName := fmt.Sprintf("%v_%v_%v_%v_%v",
		p.ProjectId,
		p.DeploymentId,
		p.Branch,
		sha,
		p.DeployedAt.UnixMilli())
	projectNameSha := sha256.Sum256([]byte(projectName))

	return fmt.Sprintf("%.*x", 12, projectNameSha)
}
