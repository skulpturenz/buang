package enumsenv

import "fmt"

type Environment int

const (
	Production Environment = iota
	Development
	Test
)

func (env Environment) String() string {
	return []string{"production", "development", "test"}[env]
}

func Parse(s string) (Environment, error) {
	switch s {
	case "production":
		return Production, nil
	case "development":
		return Development, nil
	case "test":
		return Test, nil
	}

	return Development, fmt.Errorf("unrecognized db type: %s", s)
}
