package tools

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"skulpture/buang/components/projects"
	workersinterfaces "skulpture/buang/workers/interfaces"
)

var deleteProjectTool = Tool{
	Name:        "delete_project",
	Description: "Delete a project and tear down all its active deployments",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"projectId": map[string]any{
				"type":        "number",
				"description": "Project ID",
			},
		},
		"required": []string{"projectId"},
	},
	Handler: deleteProject,
}

func deleteProject(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
	req, err := validateArgs[DeleteProjectArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	d := deployments.FindActiveDeploymentsByProjectParams{
		ProjectID: req.ProjectId,
	}

	activeDeployments, err := d.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	for _, dply := range activeDeployments.Deployments {
		wp := workersinterfaces.BuangDeploymentParams{
			ProjectId:    req.ProjectId,
			DeploymentId: dply.GetId(),
			Block:        true,
		}

		err = s.Workflows.BuangDeployment(ctx, wp)
		if err != nil {
			return "", err
		}
	}

	p := projects.DeleteProjectParams{
		Id: req.ProjectId,
	}

	_, err = p.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	return "ok", nil
}