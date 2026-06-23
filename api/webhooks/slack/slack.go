package slack

func RenderFailedDeployment(logs string) string {
	return "Deployment Failed\n```\n" + logs + "\n```"
}

func RenderSuccessfulDeployment(url string) string {
	return "Deployment Successful\nDeployed to: " + url
}