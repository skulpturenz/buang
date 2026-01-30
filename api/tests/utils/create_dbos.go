package testutils

import (
	"context"
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
)

func CreateDbos(ctx context.Context) (*DurableExecutorConfiguration, func(context.Context) error, error) {
	pg, cleanup, err := CreatePg(ctx)
	if err != nil {
		return nil, nil, err
	}

	res := DurableExecutorConfiguration{
		DurableExecutor:    enumsdurableexecutors.Dbos,
		DbType:             enumsdbtypes.Pg,
		DbConnectionString: pg.ConnectionString,
		DbContainer:        &pg.Container,
	}

	return &res, cleanup, nil
}
