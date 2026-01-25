package testutils

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type CreatePgResult struct {
	DatabaseName     string
	Username         string
	Password         string
	ConnectionString string
}

func CreatePg(ctx context.Context) (*CreatePgResult, func(context.Context) error, error) {
	dbName := "buang"
	user := "buang"
	pw := "buang"

	connectionString := fmt.Sprintf("postgresql://%v:%v@localhost/%v?sslmode=disable", user, pw, dbName)

	pg, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(pw),
		postgres.BasicWaitStrategies())
	if err != nil {
		return nil, nil, err
	}

	cleanup := func(ctx context.Context) error {
		if err := testcontainers.TerminateContainer(pg); err != nil {
			return err
		}

		return nil
	}

	res := CreatePgResult{
		DatabaseName:     dbName,
		Username:         user,
		Password:         pw,
		ConnectionString: connectionString,
	}

	return &res, cleanup, nil
}
