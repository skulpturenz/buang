package webhooks

import (
	"fmt"
	"io"

	enumswebhooktype "skulpture/buang/enums/webhook_type"
	"skulpture/buang/webhooks/discord"
	"skulpture/buang/webhooks/slack"
)

func NewSuccessfulDeployment(webhookType enumswebhooktype.WebhookType, webhookURL string, deploymentURL string) (io.Writer, fmt.Stringer) {
	switch webhookType {
	case enumswebhooktype.Slack:
		return slack.SuccessfulDeploymentDetails{
			WebhookURL:     webhookURL,
			DeploymentURL: deploymentURL,
		}, nil
	case enumswebhooktype.Discord:
		return discord.SuccessfulDeploymentDetails{
			WebhookURL:     webhookURL,
			DeploymentURL: deploymentURL,
		}, nil
	default:
		return nil, nil
	}
}

func NewFailedDeployment(webhookType enumswebhooktype.WebhookType, webhookURL string, deploymentLogs string) (io.Writer, fmt.Stringer) {
	switch webhookType {
	case enumswebhooktype.Slack:
		return slack.FailedDeploymentDetails{
			WebhookURL:     webhookURL,
			DeploymentLogs: deploymentLogs,
		}, nil
	case enumswebhooktype.Discord:
		return discord.FailedDeploymentDetails{
			WebhookURL:     webhookURL,
			DeploymentLogs: deploymentLogs,
		}, nil
	default:
		return nil, nil
	}
}