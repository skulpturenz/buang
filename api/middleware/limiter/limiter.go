package limiter

import (
	"context"
	enumsenv "skulpture/buang/enums/env"
	"time"

	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

type LimiterConfig struct {
	Env      enumsenv.Environment
	Tokens   uint64
	Interval time.Duration
	store    limiter.Store
	limiter  *httplimit.Middleware
}

func (c LimiterConfig) New(ctx context.Context) (*httplimit.Middleware, error) {
	if c.Env != enumsenv.Production {
		noopStore, err := noopstore.New()
		if err != nil {
			return nil, err
		}

		c.store = noopStore
	} else {
		memoryStore, err := memorystore.New(&memorystore.Config{
			Tokens:   c.Tokens,
			Interval: c.Interval,
		})
		if err != nil {
			return nil, err
		}

		c.store = memoryStore
	}

	limiter, err := httplimit.NewMiddleware(c.store, httplimit.IPKeyFunc("X-Forwarded-For"))
	if err != nil {
		return nil, err
	}

	return limiter, nil
}
