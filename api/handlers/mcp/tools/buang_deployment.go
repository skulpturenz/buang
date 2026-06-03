package tools

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	workersinterfaces "skulpture/buang/workers/interfaces"
)

var buangDeploymentTool = Tool{
	Name:        "buang_deployment",
	Description: "Spin down a preview deployment. Always waits for the teardown to complete before returning.",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"projectId": map[string]any{
				"type":        "number",
				"description": "Project ID",
			},
			"deploymentId": map[string]any{
				"type":        "number",
				"description": "Deployment ID",
			},
		},
		"required": []string{"projectId", "deploymentId"},
	},
	Handler: buangDeployment,
}

func buangDeployment(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
	req, err := validateArgs[BuangDeploymentArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	wp := workersinterfaces.BuangDeploymentParams{
		ProjectId:    req.ProjectId,
		DeploymentId: req.DeploymentId,
		Block:        true,
	}

	err = s.Workflows.BuangDeployment(ctx, wp)
	if err != nil {
		return "", err
	}

	return "ok", nil
}