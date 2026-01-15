package enums

type Environment int

const (
	Production Environment = iota
	Development
	Test
)

func (env Environment) String() string {
	return []string{"production", "development", "test"}[env]
}
