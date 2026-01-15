package wrappers

import (
	sqlite_models "skulpture/buang/db/sqlite/out"
)

type DeploymentLog sqlite_models.DeploymentLog

func (dl DeploymentLog) GetId() int64 {
	return dl.ID
}

func (dl DeploymentLog) GetDeploymentId() int64 {
	return dl.DeploymentID
}

func (dl DeploymentLog) GetLog() *string {
	return dl.Log
}

func (dl DeploymentLog) Unwrap() sqlite_models.DeploymentLog {
	return sqlite_models.DeploymentLog(dl)
}
