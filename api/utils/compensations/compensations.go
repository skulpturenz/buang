package compensations

import (
	"context"
	"slices"
	"sync"
)

type compensations struct {
	compensations []func(context.Context)
	once          sync.Once
	ch            chan func(context.Context)
}

func New() compensations {
	return compensations{ch: make(chan func(context.Context), 1)}
}

func (c *compensations) AddCompensation(f func(context.Context)) *compensations {
	c.ch <- f

	fn := <-c.ch
	c.compensations = append(c.compensations, fn)

	return c
}

func (c *compensations) Compensate(ctx context.Context) {
	c.once.Do(func() {
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
