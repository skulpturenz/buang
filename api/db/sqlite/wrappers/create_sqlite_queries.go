package wrappers

import (
	"context"
	"database/sql"

	"skulpture/buang/db/interfaces"
	schema "skulpture/buang/db/sqlite/schema"
	sqlite "skulpture/buang/db/sqlite/out"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type SqliteConfig struct {
	ConnectionString string
}

func (c SqliteConfig) New(_ context.Context) (interfaces.Queries, func(ctx context.Context), error) {
	db, err := sql.Open("sqlite3", c.ConnectionString)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func(ctx context.Context) {
		db.Close()
	}

	d, err := iofs.New(schema.Files, ".")
	if err != nil {
		return nil, cleanup, err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, cleanup, err
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, cleanup, err
	}

	queries := sqlite.New(db)

	return interfaces.Queries(Queries(*queries)), cleanup, nil
}