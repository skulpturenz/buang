package webhooks

import (
	enumswebhooktype "skulpture/buang/enums/webhook_type"
	"skulpture/buang/webhooks/discord"
	"skulpture/buang/webhooks/slack"
)

type RenderResult struct {
	Body string
}

func RenderFailedDeployment(webhookType enumswebhooktype.WebhookType, logs string) RenderResult {
	switch webhookType {
	case enumswebhooktype.Slack:
		return RenderResult{Body: slack.RenderFailedDeployment(logs)}
	case enumswebhooktype.Discord:
		return RenderResult{Body: discord.RenderFailedDeployment(logs)}
	default:
		return RenderResult{Body: ""}
	}
}

func RenderSuccessfulDeployment(webhookType enumswebhooktype.WebhookType, url string) RenderResult {
	switch webhookType {
	case enumswebhooktype.Slack:
		return RenderResult{Body: slack.RenderSuccessfulDeployment(url)}
	case enumswebhooktype.Discord:
		return RenderResult{Body: discord.RenderSuccessfulDeployment(url)}
	default:
		return RenderResult{Body: ""}
	}
}