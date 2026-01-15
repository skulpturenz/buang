package wrappers

import (
	sqlite_models "skulpture/buang/db/sqlite/out"
	"time"
)

type Project sqlite_models.Project

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

func (p Project) Unwrap() sqlite_models.Project {
	return sqlite_models.Project(p)
}
