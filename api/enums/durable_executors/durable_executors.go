package enumsdurableexecutors

import "fmt"

type DurableExecutor int

const (
	Temporal DurableExecutor = iota
	Dbos
)

func (env DurableExecutor) String() string {
	return []string{"temporal", "dbos"}[env]
}

func Parse(s string) (DurableExecutor, error) {
	switch s {
	case "temporal":
		return Temporal, nil
	case "dbos":
		return Dbos, nil
	}

	return Temporal, fmt.Errorf("unrecognized deployment status: %s", s)
}
