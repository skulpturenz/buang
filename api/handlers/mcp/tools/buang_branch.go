package tools

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	workersinterfaces "skulpture/buang/workers/interfaces"
)

var buangBranchTool = Tool{
	Name:        "buang_branch",
	Description: "Spin down all deployments for a branch within a project",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"projectId": map[string]any{
				"type":        "number",
				"description": "Project ID",
			},
			"branch": map[string]any{
				"type":        "string",
				"description": "Branch name to spin down",
			},
		},
		"required": []string{"projectId", "branch"},
	},
	Handler: buangBranch,
}

func buangBranch(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
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

	wp := workersinterfaces.BuangBranchParams{
		ProjectId: int64(projectIdFloat),
		Branch:    branch,
	}

	err := s.Workflows.BuangBranch(ctx, wp)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
