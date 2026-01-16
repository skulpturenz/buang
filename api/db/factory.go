package db

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgwrappers "skulpture/buang/db/pg/wrappers"
	sqlitewrappers "skulpture/buang/db/sqlite/wrappers"
	dbtypes "skulpture/buang/enums/db_types"
)

type DbConfig struct {
	Type             dbtypes.DbType
	ConnectionString string
}

func (c DbConfig) New(ctx context.Context) (interfaces.Queries, func(ctx context.Context), error) {
	if c.Type == dbtypes.Pg {
		cfg := pgwrappers.PgConfig{
			ConnectionString: c.ConnectionString,
		}

		return cfg.New(ctx)
	}

	if c.Type == dbtypes.Sqlite {
		cfg := sqlitewrappers.SqliteConfig{
			ConnectionString: c.ConnectionString,
		}

		return cfg.New(ctx)
	}

	return nil, nil, nil
}
