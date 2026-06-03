package tools

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"
)

var createProjectTool = Tool{
	Name:        "create_project",
	Description: "Create a new project for a git repository",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"repository": map[string]any{
				"type":        "string",
				"description": "Git repository URL",
			},
			"composePath": map[string]any{
				"type":        "string",
				"description": "Path to the docker-compose file within the repository",
			},
			"requiresAuthn": map[string]any{
				"type":        "boolean",
				"description": "Whether the deployed service requires basic authentication",
			},
			"username": map[string]any{
				"type":        "string",
				"description": "Basic auth username (required if requiresAuthn is true)",
			},
			"password": map[string]any{
				"type":        "string",
				"description": "Basic auth password (required if requiresAuthn is true)",
			},
		},
		"required": []string{"repository", "composePath"},
	},
	Handler: createProject,
}

func createProject(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
	req, err := validateArgs[CreateProjectArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	p := projects.CreateProjectParams{
		Repository:    req.Repository,
		ComposePath:   req.ComposePath,
		RequiresAuthn: req.RequiresAuthn,
		Username:      req.Username,
		Password:      req.Password,
	}

	res, err := p.Exec(ctx, &s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", res.Id), nil
}