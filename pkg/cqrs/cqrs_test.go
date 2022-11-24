package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"

	"github.com/stretchr/testify/require"
)

func TestChErrMw(t *testing.T) {
	t.Run(`Given a ChErrMw middleware, when the wrapped command handler returns an error, then it's logged`, func(t *testing.T) {
		var (
			logger    = &LoggerMock{}
			randomErr = errors.New("")
			ch        = &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]cqrs.Event, error) {
					return nil, randomErr
				},
			}
		)
		_, err := cqrs.ChErrMw(logger)(ch).Handle(context.Background(), &CommandMock{})
		require.ErrorIs(t, err, randomErr)
		require.Len(t, ch.HandleCalls(), 1)
		require.Len(t, logger.PrintfCalls(), 1)
	})
}

func TestQhErrMw(t *testing.T) {
	t.Run(`Given a ChErrMw middleware, when the wrapped command handler returns an error, then it's logged`, func(t *testing.T) {
		var (
			logger    = &LoggerMock{}
			randomErr = errors.New("")
			qh        = &QueryHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Query) (cqrs.QueryResult, error) {
					return nil, randomErr
				},
			}
		)
		_, err := cqrs.QhErrMw(logger)(qh).Handle(context.Background(), &QueryMock{})
		require.ErrorIs(t, err, randomErr)
		require.Len(t, qh.HandleCalls(), 1)
		require.Len(t, logger.PrintfCalls(), 1)
	})
}
