package wrappers

import (
	"encoding/json"
	"skulpture/buang/db/interfaces"
	pg_models "skulpture/buang/db/pg/out"
	"time"
)

type Deployment pg_models.Deployment

var _ interfaces.Deployment = (*Deployment)(nil)

func (d Deployment) GetId() int64 {
	return d.ID
}

func (d Deployment) GetProjectId() int64 {
	return d.ProjectID
}

func (d Deployment) GetUrl() *string {
	return d.Url
}

func (d Deployment) GetStatus() int16 {
	return d.Status
}

func (d Deployment) GetSha() string {
	return d.Sha
}

func (d Deployment) GetDeployedAt() *time.Time {
	return d.DeployedAt
}

func (d Deployment) GetClonePath() *string {
	return d.ClonePath
}

func (d Deployment) GetServiceEntrypoint() string {
	return d.ServiceEntrypoint
}

func (d Deployment) GetBranch() string {
	return d.Branch
}

func (d Deployment) GetEnvVars() (map[string]any, error) {
	env := d.EnvVars

	var res map[string]any
	err := json.Unmarshal(env, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d Deployment) Unwrap() pg_models.Deployment {
	return pg_models.Deployment(d)
}
