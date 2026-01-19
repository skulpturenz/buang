package wrappers

import (
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

type Queries pgmodels.Queries

var _ interfaces.Queries = (*Queries)(nil)
