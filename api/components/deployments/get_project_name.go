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
	PRECISION := 6
	sha := fmt.Sprintf("%.*s", 8, p.Sha)
	projectName := fmt.Sprintf("%v_%v_%v_%v_%v",
		p.ProjectId,
		p.DeploymentId,
		p.Branch,
		sha,
		p.DeployedAt.UnixMilli())
	projectNameSha := sha256.Sum256([]byte(projectName))

	// 6 * 2 = 12 chars. each hex digit is 2 chars
	return fmt.Sprintf("%.*x", PRECISION, projectNameSha)
}
