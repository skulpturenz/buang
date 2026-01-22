package wrappers

import (
	"encoding/json"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"time"
)

type DiagnosticLog sqlitemodels.DiagnosticLog

var _ interfaces.DiagnosticLog = (*DiagnosticLog)(nil)

func (d DiagnosticLog) GetId() int64 {
	return d.ID
}

func (d DiagnosticLog) GetType() (enumsdiagnosticlogtype.DiagnosticLogType, error) {
	t, err := enumsdiagnosticlogtype.Parse(d.Type)

	return t, err
}

func (d DiagnosticLog) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d DiagnosticLog) GetLog() (map[string]any, error) {
	l := d.Log

	var res map[string]any
	err := json.Unmarshal(l, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
