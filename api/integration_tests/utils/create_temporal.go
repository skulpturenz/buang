package testutils

import (
	"context"
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateSqliteTemporal(ctx context.Context) (*DurableExecutorConfiguration, func(context.Context) error, error) {
	sqlite := CreateSqlite(ctx)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "temporalio/temporal",
			ExposedPorts: []string{"7233/tcp", "8233/tcp"},
			Cmd:          []string{"server", "start-dev", "--ip", "0.0.0.0"},
			WaitingFor:   wait.ForLog("0.0.0.0:7233").WithStartupTimeout(120 * time.Second).WithPollInterval(100 * time.Millisecond),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}

	res := DurableExecutorConfiguration{
		DurableExecutor:          enumsdurableexecutors.Temporal,
		DbType:                   enumsdbtypes.Sqlite,
		DbConnectionString:       sqlite.ConnectionString,
		DurableExecutorContainer: &container,
	}

	cleanup := func(ctx context.Context) error {
		if err := container.Terminate(ctx); err != nil {
			return err
		}

		return nil
	}

	return &res, cleanup, nil
}
