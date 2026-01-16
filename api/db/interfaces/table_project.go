package interfaces

import (
	"time"
)

type Project interface {
	GetId() int64
	GetRepository() string
	GetRequiresAuthn() bool
	GetUsername() *string
	GetPassword() *string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetComposePath() string
}
