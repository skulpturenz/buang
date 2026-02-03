package wrappers

import (
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

type Queries sqlitemodels.Queries

var _ interfaces.Queries = (*Queries)(nil)
