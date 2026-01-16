package enumsdeploymentstatus

import "fmt"

type DeploymentStatus int

const (
	New DeploymentStatus = iota
	Deploying
	Deployed
	Buang
	Error
	Cancelled
)

func (env DeploymentStatus) String() string {
	return []string{"new", "deploying", "deployed", "buang", "error", "cancelled"}[env]
}

func Parse(s string) (DeploymentStatus, error) {
	switch s {
	case "new":
		return New, nil
	case "deploying":
		return Deploying, nil
	case "deployed":
		return Deployed, nil
	case "buang":
		return Buang, nil
	case "error":
		return Error, nil
	case "cancelled":
		return Cancelled, nil
	}

	return Error, fmt.Errorf("unrecognized deployment status: %s", s)
}
