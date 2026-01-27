package diagnostics

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type DockerStatsRequest struct {
}

type ContainerStats struct {
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	ImageID     string      `json:"imageId"`
	Status      string      `json:"status"`
	State       string      `json:"state"`
	CpuStats    CpuStats    `json:"cpuStats"`
	MemoryStats MemoryStats `json:"memoryStats"`
}

type CpuStats struct {
	UsagePercent float64 `json:"usagePercent"`
}

type MemoryStats struct {
	UsagePercent float64 `json:"usagePercent"`
	Usage        float64 `json:"usageMb"`
	Limit        float64 `json:"limitMb"`
}

type DockerStatsResult = map[string]ContainerStats // orderedmap.OrderedMap[string, ContainerStats]

// @summary		Docker stats
// @description	Not the most accurate readings but it should tell which services need attention
// @tags			api.v1, diagnostics
// @security		ApiKeyAuth
// @success		200	{object}	DockerStatsResult
// @failure		401
// @failure		500	{object}	string
// @router			/diagnostics/stats [get]
func DockerStats(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := docker.StatsParams{}

		res, err := p.Stats(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "docker stats", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := orderedmap.New[string, ContainerStats]()
		for pair := res.Oldest(); pair != nil; pair = pair.Next() {
			ret.Set(pair.Key, ContainerStats{
				Name:        pair.Value.Name,
				Status:      pair.Value.Status,
				State:       pair.Value.State,
				Image:       pair.Value.Image,
				ImageID:     pair.Value.ImageID,
				CpuStats:    CpuStats{UsagePercent: pair.Value.CpuStats.UsagePercent},
				MemoryStats: MemoryStats{Usage: pair.Value.MemoryStats.Usage, UsagePercent: pair.Value.MemoryStats.UsagePercent, Limit: pair.Value.MemoryStats.Limit},
			})
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "docker stats", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
