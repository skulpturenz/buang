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
	sqlite, sqliteCleanup := CreateSqlite(ctx)

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
		sqliteCleanup(ctx)

		if err := container.Terminate(ctx); err != nil {
			return err
		}

		return nil
	}

	return &res, cleanup, nil
}

func CreatePgTemporal(ctx context.Context) (*DurableExecutorConfiguration, func(context.Context) error, error) {
	pg, pgCleanup, err := CreatePg(ctx)
	if err != nil {
		return nil, nil, err
	}

	postgresSeeds, err := pg.Container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	pgPort, err := pg.Container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "temporalio/temporal",
			ExposedPorts: []string{"7233/tcp", "8233/tcp"},
			Cmd:          []string{"server", "start-dev", "--ip", "0.0.0.0"},
			WaitingFor:   wait.ForLog("0.0.0.0:7233").WithStartupTimeout(120 * time.Second).WithPollInterval(100 * time.Millisecond),
			Env: map[string]string{
				"POSTGRES_SEEDS":    postgresSeeds,
				"POSTGRES_USER":     pg.Username,
				"POSTGRES_PWD":      pg.Password,
				"DBNAME":            pg.DatabaseName,
				"VISIBILITY_DBNAME": pg.DatabaseName,
				"DB_PORT":           pgPort.Port(),
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}

	res := DurableExecutorConfiguration{
		DurableExecutor:          enumsdurableexecutors.Temporal,
		DbType:                   enumsdbtypes.Pg,
		DbConnectionString:       pg.ConnectionString,
		DurableExecutorContainer: &container,
	}

	cleanup := func(ctx context.Context) error {
		pgCleanup(ctx)

		if err := container.Terminate(ctx); err != nil {
			return err
		}

		return nil
	}

	return &res, cleanup, nil
}
