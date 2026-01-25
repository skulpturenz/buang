package testutils

import "context"

type CreateSqliteResult struct {
	ConnectionString string
}

func CreateSqlite(ctx context.Context) CreateSqliteResult {
	return CreateSqliteResult{ConnectionString: "file:test.db?_foreign_keys=true&mode=memory"}
}
