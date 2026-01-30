package testutils

import (
	"context"
	"fmt"
	"os"
	"time"
)

type CreateSqliteResult struct {
	ConnectionString string
}

func CreateSqlite(ctx context.Context) (CreateSqliteResult, func(context.Context) error) {
	filePath := fmt.Sprintf("/tmp/test_%v.db", time.Now().Nanosecond())
	connectionString := fmt.Sprintf("file:%v?_foreign_keys=true", filePath)

	cleanup := func(ctx context.Context) error {
		err := os.RemoveAll(filePath)
		if err != nil {
			return err
		}

		return nil
	}

	return CreateSqliteResult{ConnectionString: connectionString}, cleanup
}
