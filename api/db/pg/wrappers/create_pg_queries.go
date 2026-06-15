package wrappers

import (
	"context"

	"skulpture/buang/db/interfaces"
	pg "skulpture/buang/db/pg/out"
	schema "skulpture/buang/db/pg/schema"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type PgConfig struct {
	ConnectionString string
}

func (c PgConfig) New(ctx context.Context) (interfaces.Queries, func(ctx context.Context), error) {
	pool, err := pgxpool.New(ctx, c.ConnectionString)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func(ctx context.Context) {
		pool.Close()
	}

	d, err := iofs.New(schema.Files, ".")
	if err != nil {
		return nil, cleanup, err
	}

	driver, err := pgx.WithInstance(stdlib.OpenDBFromPool(pool), &pgx.Config{})
	if err != nil {
		return nil, cleanup, err
	}
	defer driver.Close()

	m, err := migrate.NewWithInstance("iofs", d, "pg", driver)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, cleanup, err
	}

	queries := pg.New(pool)

	return interfaces.Queries(Queries(*queries)), cleanup, nil
}