package wrappers

import (
	pg_models "skulpture/buang/db/pg/out"
)

type DeploymentLog pg_models.DeploymentLog

func (dl DeploymentLog) GetId() int64 {
	return dl.ID
}

func (dl DeploymentLog) GetDeploymentId() int64 {
	return dl.DeploymentID
}

func (dl DeploymentLog) GetLog() *string {
	return dl.Log
}

func (dl DeploymentLog) Unwrap() pg_models.DeploymentLog {
	return pg_models.DeploymentLog(dl)
}
