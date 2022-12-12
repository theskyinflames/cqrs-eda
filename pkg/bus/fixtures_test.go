package bus_test

import (
	"context"
	"time"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
)

func handlerFixture(result interface{}, err error) bus.Handler {
	return func(_ context.Context, _ bus.Dispatchable) (interface{}, error) {
		return result, err
	}
}

func waitForCtxCancelledHandlerFixture(result interface{}, err error) bus.Handler {
	return func(ctx context.Context, _ bus.Dispatchable) (interface{}, error) {
		for {
			select {
			case <-ctx.Done():
				return nil, err
			}
		}
	}
}

func randomTimeHandlerFixture(result interface{}, duration time.Duration, err error) bus.Handler {
	return func(_ context.Context, d bus.Dispatchable) (interface{}, error) {
		// Simulate handler execution time
		time.Sleep(duration)
		return result, err
	}
}
