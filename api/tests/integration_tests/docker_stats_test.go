package integrationtests

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"skulpture/buang/handlers/diagnostics"
	testutils "skulpture/buang/tests/utils"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestDockerStats(t *testing.T) {
	t.Parallel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(t.Context())
	require.NoError(t, err)
	defer dbosCleanup(t.Context())

	temporalSqliteConfig, temporalSqliteCleanup, err := testutils.CreateSqliteTemporal(t.Context())
	require.NoError(t, err)
	defer temporalSqliteCleanup(t.Context())

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(t.Context())
	require.NoError(t, err)
	defer temporalPgCleanup(t.Context())

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalSqliteConfig,
		"TemporalPg":     *temporalPgConfig,
	}

	retry := testutils.NewRetry(3, 500*time.Millisecond)

	for k, v := range scenarios {
		retry.Retry(t, func(t *testing.T) {
			t.Run(k, func(t *testing.T) {
				dockerStats(t, v)
			})
		})
	}
}

func dockerStats(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

	testApp, cleanup := testutils.Setup(ctx, config, nil)
	defer cleanup(ctx)

	docker := testApp.GetHttpApplication().Services.Docker

	reader, err := docker.ImagePull(ctx, "alpine", image.PullOptions{})
	require.NoError(t, err)
	io.Copy(io.Discard, reader)

	busyContainer, err := docker.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"sh", "-c", "while true; do :; done"},
	}, nil, nil, nil, "")
	require.NoError(t, err)

	idleContainer, err := docker.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"sleep", "300"},
	}, nil, nil, nil, "")
	require.NoError(t, err)

	err = docker.ContainerStart(ctx, busyContainer.ID, container.StartOptions{})
	require.NoError(t, err)
	defer docker.ContainerRemove(ctx, busyContainer.ID, container.RemoveOptions{Force: true})

	err = docker.ContainerStart(ctx, idleContainer.ID, container.StartOptions{})
	require.NoError(t, err)
	defer docker.ContainerRemove(ctx, idleContainer.ID, container.RemoveOptions{Force: true})

	time.Sleep(5 * time.Second)

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)
	statsUrl := fmt.Sprintf("%v/diagnostics/stats", baseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, statsUrl, nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var stats orderedmap.OrderedMap[string, diagnostics.ContainerStats]
	err = json.NewDecoder(res.Body).Decode(&stats)
	require.NoError(t, err)

	require.NotEmpty(t, stats)

	require.Equal(t, stats.Oldest().Key, busyContainer.ID, fmt.Sprintf("busy container is %v", stats.Oldest().Value.Name))
	require.Equal(t, stats.Newest().Key, idleContainer.ID, fmt.Sprintf("idle container is %v", stats.Newest().Value.Name))

	unsorted := []diagnostics.ContainerStats{} // assume desc
	for pair := stats.Oldest(); pair != nil; pair = pair.Next() {
		unsorted = append(unsorted, pair.Value)
	}

	sorted := unsorted
	sort.Slice(sorted, func(x int, y int) bool {
		return cmp.Or(
			cmp.Compare(sorted[x].CpuStats.UsagePercent, sorted[y].CpuStats.UsagePercent),
			cmp.Compare(sorted[x].MemoryStats.UsagePercent, sorted[y].MemoryStats.UsagePercent),
		) < 0 // asc
	})

	slices.Reverse(unsorted)
	require.Equal(t, unsorted, sorted)
}
