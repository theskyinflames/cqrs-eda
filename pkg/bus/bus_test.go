package bus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"

	"github.com/stretchr/testify/require"
)

func TestBusDispatch(t *testing.T) {
	var (
		randomErr = errors.New("")
		response  = "a response"
	)
	tests := []struct {
		name            string
		handlerName     string
		handler         bus.Handler
		dispatchable    bus.Dispatchable
		expected        any
		expectedErrFunc func(t *testing.T, err error)
	}{
		{
			name:         `Given an unknown handler, when it's called, then an error is returned`,
			handlerName:  "unknown",
			handler:      handlerFixture(nil, nil),
			dispatchable: &DispatchableMock{},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, bus.ErrNotDispatchable)
			},
		},
		{
			name: `Given a dispatchable that makes its handler returns an error, 
				when it's called, then an error is returned`,
			handlerName: "h",
			handler:     handlerFixture(nil, randomErr),
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
			name:        `Given a dispatchable, when it's called, then an error is returned`,
			handlerName: "h",
			handler:     handlerFixture(response, nil),
			dispatchable: &DispatchableMock{
				NameFunc: func() string {
					return "h"
				},
			},
			expected: response,
		},
	}

	for _, tt := range tests {
		b := bus.New()
		b.Register(tt.handlerName, tt.handler)
		response, err := b.Dispatch(context.Background(), tt.dispatchable)
		require.Equal(t, tt.expectedErrFunc == nil, err == nil)
		if err != nil {
			tt.expectedErrFunc(t, err)
			continue
		}
		require.Equal(t, tt.expected, response)
	}
}
