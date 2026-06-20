package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	"time"
)

var infoTool = Tool{
	Name:        "get_info",
	Description: "Get service version, timezone, and uptime information",
	InputSchema: map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	},
	Handler: getInfo,
}

func getInfo(_ context.Context, _ app.ApplicationServices, args map[string]any) (string, error) {
	_, err := validateArgs[InfoArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	zoneName, _ := time.Now().Zone()

	info := map[string]string{
		"version":  constantsenvs.BUANG_VERSION,
		"timezone": zoneName,
		"uptime":   constantsenvs.Uptime().String(),
	}

	b, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	return string(b), nil
}