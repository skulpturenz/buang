package interfaces

import (
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"time"
)

type DiagnosticLog interface {
	GetId() int64
	GetType() (enumsdiagnosticlogtype.DiagnosticLogType, error)
	GetCreatedAt() time.Time
	GetLog() (map[string]any, error)
}
