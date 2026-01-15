package wrappers

import (
	pg_models "skulpture/buang/db/pg/out"
)

type Deployment pg_models.Deployment

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

func (d Deployment) Unwrap() pg_models.Deployment {
	return pg_models.Deployment(d)
}
