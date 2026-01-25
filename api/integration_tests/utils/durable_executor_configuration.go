package testutils

import (
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
)

type DurableExecutorConfiguration struct {
	DurableExecutor    enumsdurableexecutors.DurableExecutor
	DbType             enumsdbtypes.DbType
	DbConnectionString string
}
