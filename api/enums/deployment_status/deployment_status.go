package enumsdeploymentstatus

import "fmt"

type DeploymentStatus int

const (
	New DeploymentStatus = iota
	Deployed
	Buang
	Error
)

func (env DeploymentStatus) String() string {
	return []string{"new", "deployed", "buang", "error"}[env]
}

func Parse(s string) (DeploymentStatus, error) {
	switch s {
	case "new":
		return New, nil
	case "deployed":
		return Deployed, nil
	case "buang":
		return Buang, nil
	case "error":
		return Error, nil
	}

	return Error, fmt.Errorf("unrecognized deployment status: %s", s)
}
