package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SuccessfulDeploymentDetails struct {
	WebhookURL string
	DeploymentURL string
}

func (d SuccessfulDeploymentDetails) String() string {
	m := map[string]any{"text": d.formatMessage()}
	jsonBytes, _ := json.Marshal(m)
	return string(jsonBytes)
}

func (d SuccessfulDeploymentDetails) Write(p []byte) (n int, err error) {
	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return len(p), nil
}

func (d SuccessfulDeploymentDetails) formatMessage() string {
	return "**Deployment Successful**\nDeployed to: " + d.DeploymentURL
}

type FailedDeploymentDetails struct {
	WebhookURL string
	DeploymentLogs string
}

func (d FailedDeploymentDetails) String() string {
	m := map[string]any{"text": d.formatMessage()}
	jsonBytes, _ := json.Marshal(m)
	return string(jsonBytes)
}

func (d FailedDeploymentDetails) Write(p []byte) (n int, err error) {
	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return len(p), nil
}

func (d FailedDeploymentDetails) formatMessage() string {
	return "**Deployment Failed**\n```\n" + d.DeploymentLogs + "\n```"
}