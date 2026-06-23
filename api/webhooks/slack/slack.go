package slack

func RenderFailedDeployment(logs string) map[string]any {
	return map[string]any{"text": "Deployment Failed\n```\n" + logs + "\n```"}
}

func RenderSuccessfulDeployment(url string) map[string]any {
	return map[string]any{"text": "Deployment Successful\nDeployed to: " + url}
}