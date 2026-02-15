package deployments

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/wordgen/wordlists"
)

type GetProjectNameParams struct {
	ProjectId    int64
	DeploymentId int64
	Branch       string
	Sha          string
	DeployedAt   time.Time
}

func GetProjectName(p GetProjectNameParams) string {
	projectName := fmt.Sprintf("%v_%v_%v_%v_%v",
		p.ProjectId,
		p.DeploymentId,
		p.Branch,
		p.Sha,
		p.DeployedAt.UnixMilli())
	projectNameSha := sha256.Sum256([]byte(projectName))

	words := []string{}
	nProjectName := new(big.Int).SetBytes(projectNameSha[:])
	nWordList := big.NewInt(int64(len(wordlists.NamesMixed)))

	for range 4 {
		remainder := new(big.Int)
		nProjectName.QuoRem(nProjectName, nWordList, remainder)
		words = append(words, wordlists.NamesMixed[remainder.Int64()])
	}

	return strings.Join(words, "-")
}
