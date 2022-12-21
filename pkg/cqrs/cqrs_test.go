package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/events"

	"github.com/stretchr/testify/require"
)

func TestChErrMw(t *testing.T) {
	t.Run(`Given a ChErrMw middleware, when the wrapped command handler returns an error, then it's logged`, func(t *testing.T) {
		var (
			logger    = &LoggerMock{}
			randomErr = errors.New("")
			ch        = &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]events.Event, error) {
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

func TestCommandHandlerMultiMiddleware(t *testing.T) {
	t.Run(`Given a sequence of ch middlewares,
		when it's called, 
		then the ch is executed wrapped by all middlewares in the right order`, func(t *testing.T) {
		var (
			calls []string
			chMw1 = chTestMw("mw1", &calls)
			chMw2 = chTestMw("mw2", &calls)
			chMw3 = chTestMw("mw3", &calls)
			chMw4 = chTestMw("mw4", &calls)

			ev = &EventMock{}

			ch = &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]events.Event, error) {
					return []events.Event{ev}, nil
				},
			}
		)

		multiChMw := cqrs.CommandHandlerMultiMiddleware(chMw1, chMw2, chMw3, chMw4)

		evs, err := multiChMw(ch).Handle(context.Background(), &CommandMock{})

		require.Len(t, ch.HandleCalls(), 1)
		require.NoError(t, err)
		require.Len(t, evs, 1)
		require.Equal(t, ev, evs[0])
		require.Equal(t, []string{"mw4", "mw3", "mw2", "mw1"}, calls)
	})
}

func chTestMw(name string, calls *[]string) cqrs.CommandHandlerMiddleware {
	return func(ch cqrs.CommandHandler) cqrs.CommandHandler {
		return cqrs.CommandHandlerFunc(func(ctx context.Context, cmd cqrs.Command) ([]events.Event, error) {
			*calls = append(*calls, name)
			return ch.Handle(ctx, cmd)
		})
	}
}

func TestQueryHandlerMultiMiddleware(t *testing.T) {
	t.Run(`Given a sequence of ch middlewares,
		when it's called, 
		then the ch is executed wrapped by all middlewares in the right order`, func(t *testing.T) {
		var (
			calls []string
			qhMw1 = qhTestMw("mw1", &calls)
			qhMw2 = qhTestMw("mw2", &calls)
			qhMw3 = qhTestMw("mw3", &calls)
			qhMw4 = qhTestMw("mw4", &calls)

			queryResult = "result"

			qh = &QueryHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Query) (cqrs.QueryResult, error) {
					return queryResult, nil
				},
			}
		)

		multiQhMw := cqrs.QueryHandlerMultiMiddleware(qhMw1, qhMw2, qhMw3, qhMw4)

		qrs, err := multiQhMw(qh).Handle(context.Background(), &QueryMock{})

		require.Len(t, qh.HandleCalls(), 1)
		require.NoError(t, err)
		require.Equal(t, queryResult, qrs)
	})
}

func qhTestMw(name string, calls *[]string) cqrs.QueryHandlerMiddleware {
	return func(ch cqrs.QueryHandler) cqrs.QueryHandler {
		return cqrs.QueryHandlerFunc(func(ctx context.Context, cmd cqrs.Query) (cqrs.QueryResult, error) {
			*calls = append(*calls, name)
			return ch.Handle(ctx, cmd)
		})
	}
}

func TestChEventMw(t *testing.T) {
	t.Run(`Given a events ch middleware with a events bus,
	when it catches an error from the ch, 
	then no events are dispatched to the events bus`, func(t *testing.T) {
		var (
			eventName = "entity.changed"
			ev        = &EventMock{
				NameFunc: func() string {
					return eventName
				},
			}
			evBus = &BusMock{
				DispatchFunc: func(_ context.Context, _ bus.Dispatchable) (interface{}, error) {
					return nil, nil
				},
			}
			err = errors.New("")
			ch  = &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]events.Event, error) {
					return []events.Event{ev}, err
				},
			}
		)

		evs, gotErr := cqrs.ChEventMw(evBus)(ch).Handle(context.Background(), &CommandMock{})
		require.ErrorIs(t, err, gotErr)
		require.Len(t, ch.HandleCalls(), 1)
		require.Len(t, evBus.DispatchCalls(), 0)
		require.Len(t, evs, 1)
		require.Equal(t, ev, evs[0])
	})

	t.Run(`Given a events ch middleware with a events bus,
	when it catches events from the ch, 
	then they are dispatched to the events bus`, func(t *testing.T) {
		var (
			eventName = "entity.changed"
			ev        = &EventMock{
				NameFunc: func() string {
					return eventName
				},
			}
			evBus = &BusMock{
				DispatchFunc: func(_ context.Context, _ bus.Dispatchable) (interface{}, error) {
					return nil, nil
				},
			}
			ch = &CommandHandlerMock{
				HandleFunc: func(_ context.Context, _ cqrs.Command) ([]events.Event, error) {
					return []events.Event{ev}, nil
				},
			}
		)

		evs, err := cqrs.ChEventMw(evBus)(ch).Handle(context.Background(), &CommandMock{})
		require.NoError(t, err)
		require.Len(t, ch.HandleCalls(), 1)
		require.Len(t, evBus.DispatchCalls(), 1)
		require.Equal(t, ev, evBus.DispatchCalls()[0].Dispatchable)
		require.Len(t, evs, 1)
		require.Equal(t, ev, evs[0])
	})
}
