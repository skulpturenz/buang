package temporalworkflows

type NonRetryableError int

const (
	UhOh NonRetryableError = iota
)

func (env NonRetryableError) String() string {
	return []string{"uhoh"}[env]
}

func Parse(s string) NonRetryableError {
	switch s {
	case "uhoh":
		return UhOh
	}

	return UhOh
}
