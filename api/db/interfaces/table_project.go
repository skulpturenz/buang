package interfaces

import (
	"time"
)

type Project interface {
	GetId() int64
	GetRepository() string
	GetRequiresAuth() bool
	GetUsername() *string
	GetPassword() *string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
