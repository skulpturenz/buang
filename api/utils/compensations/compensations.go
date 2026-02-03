package compensations

import (
	"context"
	"slices"
	"sync"
)

type compensations struct {
	compensations []func(context.Context)
	once          sync.Once
	mu            sync.Mutex
}

func New() *compensations {
	return &compensations{}
}

func (c *compensations) AddCompensation(f func(context.Context)) *compensations {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.compensations = append(c.compensations, f)

	return c
}

func (c *compensations) Compensate(ctx context.Context) {
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		for _, v := range slices.Backward(c.compensations) {
			v(ctx)
		}
	})
}

func (c *compensations) CompensateAndPanic(ctx context.Context, err error) {
	c.Compensate(ctx)

	panic(err)
}

func (c *compensations) CompensateAndError(ctx context.Context, err error) error {
	c.Compensate(ctx)

	return err
}
