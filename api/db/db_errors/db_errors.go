package dberrors

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
)

func IsNoRows(x error) bool {
	return errors.Is(sql.ErrNoRows, x) || errors.Is(pgx.ErrNoRows, x)
}
