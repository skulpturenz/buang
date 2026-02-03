package enumsdbtypes

import "fmt"

type DbType int

const (
	Pg DbType = iota
	Sqlite
)

func (env DbType) String() string {
	return []string{"postgres", "sqlite"}[env]
}

func Parse(s string) (DbType, error) {
	switch s {
	case "postgres":
		return Pg, nil
	case "sqlite":
		return Sqlite, nil
	}

	return Sqlite, fmt.Errorf("unrecognized db type: %s", s)
}
