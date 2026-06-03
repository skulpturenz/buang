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
	req, err := validateArgs[BuangBranchArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	wp := workersinterfaces.BuangBranchParams{
		ProjectId: req.ProjectId,
		Branch:    req.Branch,
	}

	err = s.Workflows.BuangBranch(ctx, wp)
	if err != nil {
		return "", err
	}

	return "ok", nil
}