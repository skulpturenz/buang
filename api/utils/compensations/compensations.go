package compensations

import (
	"context"
	"slices"
)

type compensations struct {
	compensations []func(context.Context)
}

func New() compensations {
	return compensations{}
}

func (c *compensations) AddCompensation(f func(context.Context)) *compensations {
	c.compensations = append(c.compensations, f)

	return c
}

func (c compensations) Compensate(ctx context.Context) {
	for _, v := range slices.Backward(c.compensations) {
		v(ctx)
	}
}

func (c compensations) CompensateAndPanic(ctx context.Context, err error) {
	for _, v := range slices.Backward(c.compensations) {
		v(ctx)
	}

	panic(err)
}

func (c compensations) CompensateAndError(ctx context.Context, err error) error {
	for _, v := range slices.Backward(c.compensations) {
		v(ctx)
	}

	return err
}
