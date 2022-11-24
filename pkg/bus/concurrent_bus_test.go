package bus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"

	"github.com/stretchr/testify/require"
)

func TestConcurrentBus(t *testing.T) {
	t.Run(`Given a concurrent `, func(t *testing.T) {
		type durationFunc func() time.Duration
		var (
			df = durationFunc(func() time.Duration { return time.Hour })

			randomErr       = errors.New("")
			ctxCancelledErr = errors.New("ctx cancelled")
			response        = "a response"
		)
		tests := []struct {
			name             string
			timeout          durationFunc
			concurrencyLimit int
			handler          bus.Handler
			handlerName      string
			dispatchable     bus.Dispatchable
			expected         bus.Response
			expectedErrFunc  func(t *testing.T, err error)
		}{
			{
				name: `with an unknown handler, 
					when it's called,
					then an error is returned`,
				timeout:          df,
				concurrencyLimit: 1,
				handlerName:      "unknown",
				handler:          handlerFixture(nil, nil),
				dispatchable:     &DispatchableMock{},
				expectedErrFunc: func(t *testing.T, err error) {
					require.ErrorIs(t, err, bus.ErrNotDispatchable)
				},
			},
			{
				name: `with a handler that returns an error, 
					when it's called, 
					then an error is returned`,
				timeout:          df,
				concurrencyLimit: 1,
				handlerName:      "h",
				handler:          handlerFixture(nil, randomErr),
				dispatchable: &DispatchableMock{
					NameFunc: func() string {
						return "h"
					},
				},
				expectedErrFunc: func(t *testing.T, err error) {
					require.ErrorIs(t, err, randomErr)
				},
			},
			{
				name: `,when it's called with a distpatchable that takes more time than bus timeout, 
					then the dispatchable is cancelled`,
				timeout:          func() time.Duration { return time.Millisecond },
				concurrencyLimit: 1,
				handlerName:      "h",
				handler:          waitForCtxCancelledHandlerFixture(nil, ctxCancelledErr),
				dispatchable: &DispatchableMock{
					NameFunc: func() string {
						return "h"
					},
				},
				expectedErrFunc: func(t *testing.T, err error) {
					require.ErrorIs(t, err, ctxCancelledErr)
				},
			},
			{
				name:             ` when it's called, then no error is returned`,
				timeout:          df,
				concurrencyLimit: 1,
				handlerName:      "h",
				handler:          handlerFixture(response, nil),
				dispatchable: &DispatchableMock{
					NameFunc: func() string {
						return "h"
					},
				},
				expected: bus.Response{Response: response},
			},
		}

		for _, tt := range tests {
			bus := bus.NewConcurrentBus(tt.timeout(), tt.concurrencyLimit)
			bus.Register(tt.handlerName, tt.handler)
			ctx, cancel := context.WithCancel(context.Background())
			go bus.Run(ctx)
			defer cancel()

			rs := <-bus.Dispatch(ctx, tt.dispatchable)
			require.Equal(t, tt.expectedErrFunc == nil, rs.Err == nil)
			if rs.Err != nil {
				tt.expectedErrFunc(t, rs.Err)
			}
		}
	})

	t.Run(`Given a concurrent bus, `, func(t *testing.T) {
		t.Parallel()
		type durationFunc func() time.Duration
		var (
			hndNames = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
			hnds     = []bus.Handler{
				randomTimeHandlerFixture(hndNames[0], time.Millisecond*100, nil),
				randomTimeHandlerFixture(hndNames[1], time.Millisecond*50, nil),
				randomTimeHandlerFixture(hndNames[2], time.Millisecond*20, nil),
				randomTimeHandlerFixture(hndNames[3], time.Millisecond*5, nil),
				randomTimeHandlerFixture(hndNames[4], time.Millisecond*40, nil),
				randomTimeHandlerFixture(hndNames[5], time.Millisecond*1, nil),
				randomTimeHandlerFixture(hndNames[6], time.Millisecond*200, nil),
				randomTimeHandlerFixture(hndNames[7], time.Millisecond*15, nil),
				randomTimeHandlerFixture(hndNames[8], time.Millisecond*4, nil),
				randomTimeHandlerFixture(hndNames[9], time.Millisecond*100, nil),
			}
		)

		const maxConcurrentSize = 2

		t.Run(`with set of concurrent dispatchable items and a bus with a limited concurrency limit, 
			when it's called with a more dispatchable than concurrency limit, 
			then concurrency limit is not overcame`, func(t *testing.T) {
			cbus := bus.NewConcurrentBus(time.Hour, maxConcurrentSize)

			for i := range hnds {
				cbus.Register(hndNames[i], hnds[i])
			}
			ctx, cancel := context.WithCancel(context.Background())
			go cbus.Run(ctx)
			defer cancel()

			// Ensure that the maximum concurrent size is not exceeded
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						require.LessOrEqual(t, cbus.CurrentSize(), maxConcurrentSize)
					}
				}
			}()

			// Dispatch items to the bus
			rsChan := make(chan bus.Response, len(hndNames))
			for _, d := range []bus.Dispatchable{
				&DispatchableMock{NameFunc: func() string { return hndNames[0] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[1] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[2] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[3] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[4] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[5] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[6] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[7] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[8] }},
				&DispatchableMock{NameFunc: func() string { return hndNames[9] }},
			} {
				// Simulate tasks dispatching items to the bus randomly
				go func(d bus.Dispatchable) {
					rsChan <- <-(cbus.Dispatch(ctx, d))
				}(d)
			}

			var responses []bus.Response
			for {
				responses = append(responses, <-rsChan)
				if len(responses) == len(hndNames) {
					break
				}
			}

			require.Len(t, responses, len(hndNames))
		})
	})
}
