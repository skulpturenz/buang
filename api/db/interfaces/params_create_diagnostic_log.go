package interfaces

import enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"

type CreateDiagnosticLogParams struct {
	Type enumsdiagnosticlogtype.DiagnosticLogType
	Log  map[string]any
}
