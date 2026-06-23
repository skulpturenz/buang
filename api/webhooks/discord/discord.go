package discord

func RenderFailedDeployment(logs string) map[string]any {
	return map[string]any{"content": "**Deployment Failed**\n```\n" + logs + "\n```"}
}

func RenderSuccessfulDeployment(url string) map[string]any {
	return map[string]any{"content": "**Deployment Successful**\nDeployed to: " + url}
}