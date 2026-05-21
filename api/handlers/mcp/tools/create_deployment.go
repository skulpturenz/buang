package tools

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	workersinterfaces "skulpture/buang/workers/interfaces"
)

var createDeploymentTool = Tool{
	Name:        "create_deployment",
	Description: "Spin up a preview deployment for a branch. Always waits for the deployment to complete before returning.",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"projectId": map[string]any{
				"type":        "number",
				"description": "Project ID",
			},
			"branch": map[string]any{
				"type":        "string",
				"description": "Branch name to deploy",
			},
			"sha": map[string]any{
				"type":        "string",
				"description": "Git commit SHA to deploy",
			},
			"serviceEntrypoint": map[string]any{
				"type":        "string",
				"description": "Host:port of the primary service entrypoint (e.g. app:8080)",
			},
			"env": map[string]any{
				"type":        "object",
				"description": "Optional environment variables to inject into the deployment",
			},
		},
		"required": []string{"projectId", "branch", "sha", "serviceEntrypoint"},
	},
	Handler: createDeployment,
}

func createDeployment(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
	projectIdRaw, ok := args["projectId"]
	if !ok {
		return "", fmt.Errorf("projectId is required")
	}
	projectIdFloat, ok := projectIdRaw.(float64)
	if !ok {
		return "", fmt.Errorf("projectId must be a number")
	}

	branch, ok := args["branch"].(string)
	if !ok || branch == "" {
		return "", fmt.Errorf("branch is required")
	}

	sha, ok := args["sha"].(string)
	if !ok || sha == "" {
		return "", fmt.Errorf("sha is required")
	}

	serviceEntrypoint, ok := args["serviceEntrypoint"].(string)
	if !ok || serviceEntrypoint == "" {
		return "", fmt.Errorf("serviceEntrypoint is required")
	}

	p := deployments.CreateDeploymentParams{
		ProjectID:         int64(projectIdFloat),
		Branch:            branch,
		Sha:               sha,
		ServiceEntrypoint: serviceEntrypoint,
	}

	if v, ok := args["env"].(map[string]any); ok {
		p.Env = v
	}

	res, err := p.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	wp := workersinterfaces.CreateDeploymentParams{
		ProjectId:    p.ProjectID,
		DeploymentId: res.Id,
		Block:        true,
	}

	err = s.Workflows.CreateDeployment(ctx, wp)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", res.Id), nil
}
