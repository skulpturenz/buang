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
	projectIdRaw, ok := args["projectId"]
	if !ok {
		return "", fmt.Errorf("projectId is required")
	}
	projectIdFloat, ok := projectIdRaw.(float64)
	if !ok {
		return "", fmt.Errorf("projectId must be a number")
	}
	projectId := int64(projectIdFloat)

	d := deployments.FindActiveDeploymentsByProjectParams{
		ProjectID: projectId,
	}

	activeDeployments, err := d.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	for _, dply := range activeDeployments.Deployments {
		wp := workersinterfaces.BuangDeploymentParams{
			ProjectId:    projectId,
			DeploymentId: dply.GetId(),
			Block:        true,
		}

		err = s.Workflows.BuangDeployment(ctx, wp)
		if err != nil {
			return "", err
		}
	}

	p := projects.DeleteProjectParams{
		Id: projectId,
	}

	_, err = p.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
