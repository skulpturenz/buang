package diagnostics

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/docker"
)

type DockerStatsRequest struct {
}

type ContainerStats struct {
	Name        string      `json:"name"`
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

type DockerStatsResult = map[string]ContainerStats

// @summary	Docker stats
// @tags		api.v1, diagnostics
// @security	ApiKeyAuth
// @success	200	{object}	DockerStatsResult
// @failure	401
// @failure	500	{object}	string
// @router		/diagnostics/stats [get]
func DockerStats(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := docker.StatsParams{}

		res, err := p.Stats(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "docker stats", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := map[string]ContainerStats{}
		for k, v := range res {
			ret[k] = ContainerStats{
				Name:        v.Name,
				CpuStats:    CpuStats{UsagePercent: v.CpuStats.UsagePercent},
				MemoryStats: MemoryStats{Usage: v.MemoryStats.Usage, UsagePercent: v.MemoryStats.UsagePercent, Limit: v.MemoryStats.Limit},
			}
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "docker stats", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
