package enumsauthnmethod

type AuthnMethod int

const (
	ApiKey AuthnMethod = iota
)

func (env AuthnMethod) String() string {
	return []string{"api-key"}[env]
}
