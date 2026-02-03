package testutils

import (
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"

	"github.com/testcontainers/testcontainers-go"
)

type DurableExecutorConfiguration struct {
	DurableExecutor          enumsdurableexecutors.DurableExecutor
	DbType                   enumsdbtypes.DbType
	DbConnectionString       string
	DbContainer              *testcontainers.Container
	DurableExecutorContainer *testcontainers.Container
}
