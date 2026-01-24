package docker

import (
	"cmp"
	"context"
	"encoding/json"
	"skulpture/buang/app"
	"sort"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

type StatsParams struct {
	All *bool
}

type ContainerStats struct {
	id          string
	Name        string
	CpuStats    CpuStats
	MemoryStats MemoryStats
}

type CpuStats struct {
	UsagePercent float64
}

type MemoryStats struct {
	UsagePercent float64
	Usage        float64
	Limit        float64
}

type StatsResult = orderedmap.OrderedMap[string, ContainerStats]

// reference: https://docs.docker.com/reference/api/engine/version/v1.45/#tag/Container/operation/ContainerStats
func (c StatsParams) Stats(ctx context.Context, s *app.ApplicationServices) (*StatsResult, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	opts := container.ListOptions{}
	if c.All != nil {
		opts.All = *c.All
	}

	ps, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	statsChan := make(chan ContainerStats, len(ps))

	getStats := func(ctx context.Context, wg *sync.WaitGroup, id string) {
		defer wg.Done()

		s, err := cli.ContainerStats(ctx, id, true)
		if err != nil {
			return
		}
		defer s.Body.Close()

		d := json.NewDecoder(s.Body)

		var x, y container.StatsResponse
		if err := d.Decode(&x); err != nil {
			return
		}
		if err := d.Decode(&y); err != nil {
			return
		}

		// NOTE: PRECPU STATS ARE NOT AVAILABLE ON NEWER HOSTS (CGROUP V2)
		deltaCpu := float64(y.CPUStats.CPUUsage.TotalUsage) - float64(x.CPUStats.CPUUsage.TotalUsage) // ns
		deltaSysCpu := float64(y.CPUStats.SystemUsage) - float64(x.CPUStats.SystemUsage)              // ns
		numCpus := float64(x.CPUStats.OnlineCPUs)
		// (cpu_delta / system_cpu_delta) * number_cpus * 100.0
		// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
		// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
		// number_cpus = cpu_stats.online_cpus
		usagePercent := float64(0)
		if deltaSysCpu > 0 {
			usagePercent = (deltaCpu / deltaSysCpu) * numCpus * 100.0
		}

		// `memory_stats.stats.cache` is `memory.stats.total_inactive_file` on modern distributions (cgroup v2): https://docs.docker.com/reference/cli/docker/container/stats/
		memUsage := float64(x.MemoryStats.Usage) - float64(x.MemoryStats.Stats["total_inactive_file"]) // bytes
		memLimit := float64(x.MemoryStats.Limit)                                                       // bytes
		memUsagePercent := float64(0)
		// (used_memory / available_memory) * 100.0
		// used_memory = memory_stats.usage - memory_stats.stats.cache
		// available_memory = memory_stats.limit
		if memLimit > 0 {
			memUsagePercent = (memUsage / memLimit) * 100.0
		}

		select {
		case statsChan <- ContainerStats{
			id:   x.ID,
			Name: strings.Replace(x.Name, "/", "", 1),
			CpuStats: CpuStats{
				UsagePercent: usagePercent,
			},
			MemoryStats: MemoryStats{
				UsagePercent: memUsagePercent,
				Usage:        memUsage / MB,
				Limit:        memLimit / MB,
			},
		}:
		case <-ctx.Done():
			return
		}
	}

	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(statsChan)
	}()

	for _, v := range ps {
		wg.Add(1)
		go getStats(ctx, &wg, v.ID)
	}

	results := []ContainerStats{}
	for c := range statsChan {
		results = append(results, c)
	}
	sort.Slice(results, func(x int, y int) bool {
		return cmp.Or(
			cmp.Compare(results[y].CpuStats.UsagePercent, results[x].CpuStats.UsagePercent),
			cmp.Compare(results[y].MemoryStats.UsagePercent, results[x].MemoryStats.UsagePercent),
		) > 0 // desc
	})

	stats := orderedmap.New[string, ContainerStats]()
	for _, v := range results {
		stats.Set(v.id, v)
	}

	return stats, nil
}
