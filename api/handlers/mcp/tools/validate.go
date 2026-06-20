package tools

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func validateArgs[T any](args map[string]any) (T, error) {
	var empty T

	data, err := json.Marshal(args)
	if err != nil {
		return empty, fmt.Errorf("failed to marshal args: %w", err)
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return empty, fmt.Errorf("failed to unmarshal args: %w", err)
	}

	if err := validate.Struct(result); err != nil {
		return empty, err
	}

	return result, nil
}

type CreateProjectArgs struct {
	Repository    string  `json:"repository" validate:"required"`
	RequiresAuthn bool    `json:"requiresAuthn"`
	Username      *string `json:"username,omitempty" validate:"required_with=Password required_if=RequiresAuthn true"`
	Password      *string `json:"password,omitempty" validate:"required_with=Username required_if=RequiresAuthn true"`
	ComposePath   string  `json:"composePath" validate:"required"`
}

type BuangBranchArgs struct {
	ProjectId int64  `json:"projectId" validate:"required"`
	Branch    string `json:"branch" validate:"required"`
}

type DeleteProjectArgs struct {
	ProjectId int64 `json:"projectId" validate:"required"`
}

type CreateDeploymentArgs struct {
	ProjectId         int64         `json:"projectId" validate:"required"`
	Branch            string        `json:"branch" validate:"required"`
	Sha               string        `json:"sha" validate:"required"`
	ServiceEntrypoint string        `json:"serviceEntrypoint" validate:"required,hostname_port"`
	Env               map[string]any `json:"env,omitempty"`
}

type BuangDeploymentArgs struct {
	ProjectId    int64 `json:"projectId" validate:"required"`
	DeploymentId int64 `json:"deploymentId" validate:"required"`
}

type DockerStatsArgs struct {
}

type InfoArgs struct {
}