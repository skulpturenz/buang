package enumsdiagnosticlogtype

import "fmt"

type DiagnosticLogType int

const (
	DockerStats DiagnosticLogType = iota
	Panic
	Unknown
)

func (env DiagnosticLogType) String() string {
	return []string{"docker_stats", "panic"}[env]
}

func Parse(s string) (DiagnosticLogType, error) {
	switch s {
	case "docker_stats":
		return DockerStats, nil
	case "panic":
		return Panic, nil
	}

	return Unknown, fmt.Errorf("unrecognized diagnostic log type: %s", s)
}
