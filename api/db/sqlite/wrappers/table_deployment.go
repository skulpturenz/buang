package wrappers

import (
	sqlite_models "skulpture/buang/db/sqlite/out"
)

type Deployment sqlite_models.Deployment

func (d Deployment) GetId() int64 {
	return d.ID
}

func (d Deployment) GetRepositoryId() int64 {
	return d.RepositoryID
}

func (d Deployment) GetUrl() *string {
	return d.Url
}

func (d Deployment) GetStatus() int16 {
	return d.Status
}

func (d Deployment) GetSha() *string {
	return d.Sha
}

func (d Deployment) Unwrap() sqlite_models.Deployment {
	return sqlite_models.Deployment(d)
}
