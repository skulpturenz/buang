package webhooks

import (
	enumswebhooktype "skulpture/buang/enums/webhook_type"
	"skulpture/buang/webhooks/discord"
	"skulpture/buang/webhooks/slack"
)

type RenderResult map[string]any

func RenderFailedDeployment(webhookType enumswebhooktype.WebhookType, logs string) RenderResult {
	switch webhookType {
	case enumswebhooktype.Slack:
		return slack.RenderFailedDeployment(logs)
	case enumswebhooktype.Discord:
		return discord.RenderFailedDeployment(logs)
	default:
		return nil
	}
}

func RenderSuccessfulDeployment(webhookType enumswebhooktype.WebhookType, url string) RenderResult {
	switch webhookType {
	case enumswebhooktype.Slack:
		return slack.RenderSuccessfulDeployment(url)
	case enumswebhooktype.Discord:
		return discord.RenderSuccessfulDeployment(url)
	default:
		return nil
	}
}