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
	projectIdRaw, ok := args["projectId"]
	if !ok {
		return "", fmt.Errorf("projectId is required")
	}
	projectIdFloat, ok := projectIdRaw.(float64)
	if !ok {
		return "", fmt.Errorf("projectId must be a number")
	}

	deploymentIdRaw, ok := args["deploymentId"]
	if !ok {
		return "", fmt.Errorf("deploymentId is required")
	}
	deploymentIdFloat, ok := deploymentIdRaw.(float64)
	if !ok {
		return "", fmt.Errorf("deploymentId must be a number")
	}

	wp := workersinterfaces.BuangDeploymentParams{
		ProjectId:    int64(projectIdFloat),
		DeploymentId: int64(deploymentIdFloat),
		Block:        true,
	}

	err := s.Workflows.BuangDeployment(ctx, wp)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
