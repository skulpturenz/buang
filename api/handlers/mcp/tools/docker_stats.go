package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"
)

var dockerStatsTool = Tool{
	Name:        "docker_stats",
	Description: "Get CPU and memory stats for all running containers",
	InputSchema: map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	},
	Handler: dockerStats,
}

type containerStatsResult struct {
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	ImageID    string  `json:"imageId"`
	Status     string  `json:"status"`
	State      string  `json:"state"`
	CpuPercent float64 `json:"cpuUsagePercent"`
	MemPercent float64 `json:"memUsagePercent"`
	MemUsageMb float64 `json:"memUsageMb"`
	MemLimitMb float64 `json:"memLimitMb"`
}

func dockerStats(ctx context.Context, s app.ApplicationServices, args map[string]any) (string, error) {
	_, err := validateArgs[DockerStatsArgs](args)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	p := docker.StatsParams{}

	res, err := p.Stats(ctx, &s)
	if err != nil {
		return "", err
	}

	out := map[string]containerStatsResult{}
	for pair := res.Oldest(); pair != nil; pair = pair.Next() {
		v := pair.Value
		out[v.Name] = containerStatsResult{
			Name:       v.Name,
			Image:      v.Image,
			ImageID:    v.ImageID,
			Status:     v.Status,
			State:      v.State,
			CpuPercent: v.CpuStats.UsagePercent,
			MemPercent: v.MemoryStats.UsagePercent,
			MemUsageMb: v.MemoryStats.Usage,
			MemLimitMb: v.MemoryStats.Limit,
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}

	return string(b), nil
}