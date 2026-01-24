package o11y

import (
	"log/slog"
	"skulpture/buang/db/interfaces"
	"sync"

	"github.com/negrel/assert"
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

func getInstance() *diagnosticLog {
	if instance == nil {
		slog.Error("diagnostic log singleton not initialized, nothing will be captured")

		return nil
	}

	return instance
}

func (d diagnosticLog) AssertInitialized() {
	assert.True(instance != nil, "diagnostic log not initialized")
}
