package o11y

import (
	"log/slog"
	"skulpture/buang/db/interfaces"
	"sync"
)

type diagnosticLog struct {
	q *interfaces.Queries
}

var instance *diagnosticLog

var once sync.Once

func NewS(q *interfaces.Queries) *diagnosticLog {
	once.Do(func() {
		instance = &diagnosticLog{
			q: q,
		}
	})

	return instance
}

// ignoring errors here technically breaks the rules of workflows
// but its so that we can capture the information we need
func getInstance() *diagnosticLog {
	if instance == nil {
		slog.Error("diagnostic log singleton not initialized, nothing will be captured")

		return nil
	}

	return instance
}
