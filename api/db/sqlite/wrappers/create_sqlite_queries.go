package wrappers

import (
	"context"
	"database/sql"
	"skulpture/buang/db/interfaces"
	sqlite "skulpture/buang/db/sqlite/out"
	migrations "skulpture/buang/db/sqlite/schema"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
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

	s := bindata.Resource(migrations.AssetNames(), func(name string) ([]byte, error) {
		return migrations.Asset(name)
	})
	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, nil, err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, nil, err
	}

	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite", driver)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, nil, err
	}

	queries := sqlite.New(db)

	return interfaces.Queries(Queries(*queries)), cleanup, nil
}
