package wrappers

import (
	pg_models "skulpture/buang/db/pg/out"
	"time"
)

type Project pg_models.Project

func (p Project) GetId() int64 {
	return p.ID
}

func (p Project) GetRepository() string {
	return p.Repository
}

func (p Project) GetRequiresAuthn() bool {
	return p.RequiresAuthn
}

func (p Project) GetUsername() *string {
	return p.Username
}

func (p Project) GetPassword() *string {
	return p.Password
}

func (p Project) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p Project) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p Project) GetComposePath() string {
	return p.ComposePath
}

func (p Project) GetDeleted() bool {
	return p.Deleted
}

func (p Project) Unwrap() pg_models.Project {
	return pg_models.Project(p)
}
