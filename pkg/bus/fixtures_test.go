package bus_test

import (
	"context"
	"time"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
)

func handlerFixture(result any, err error) bus.Handler {
	return func(_ context.Context, _ bus.Dispatchable) (any, error) {
		return result, err
	}
}

func waitForCtxCancelledHandlerFixture(result any, err error) bus.Handler {
	return func(ctx context.Context, _ bus.Dispatchable) (any, error) {
		for {
			select {
			case <-ctx.Done():
				return nil, err
			}
		}
	}
}

func randomTimeHandlerFixture(result any, duration time.Duration, err error) bus.Handler {
	return func(_ context.Context, d bus.Dispatchable) (any, error) {
		// Simulate handler execution time
		time.Sleep(duration)
		return result, err
	}
}
