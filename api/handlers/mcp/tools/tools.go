package tools

import (
	"context"
	"skulpture/buang/app"
)

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Handler     func(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error)
}

var AllTools = []Tool{
	createProjectTool,
	buangBranchTool,
	deleteProjectTool,
	createDeploymentTool,
	buangDeploymentTool,
	dockerStatsTool,
	infoTool,
}
