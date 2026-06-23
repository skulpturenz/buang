package enumswebhooktype

import "fmt"

type WebhookType int

const (
	Slack WebhookType = iota
	Discord
)

func (env WebhookType) String() string {
	return []string{"slack", "discord"}[env]
}

func Parse(s string) (WebhookType, error) {
	switch s {
	case "slack":
		return Slack, nil
	case "discord":
		return Discord, nil
	}

	return Slack, fmt.Errorf("unrecognized webhook type: %s", s)
}
