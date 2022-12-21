package cqrs

//go:generate moq -stub -out mock_cqrs_test.go -pkg cqrs_test . Command Query CommandHandler QueryHandler Bus
//go:generate moq -stub -out mock_logger_test.go -pkg cqrs_test . Logger

import (
	"context"
	"encoding/json"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

//go:generate moq -stub -out zmock_cqrs_event_test.go -pkg cqrs_test . Event

// Logger is an interface
type Logger interface {
	Printf(format string, v ...interface{})
}

// Command is a CQRS command
type Command interface {
	Name() string
}

// CommandHandler handles a command
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) ([]events.Event, error)
}

// CommandHandlerFunc is a function that implements CommandHandler interface
type CommandHandlerFunc func(ctx context.Context, cmd Command) ([]events.Event, error)

// Handle implements the CommandHandler interface
func (chf CommandHandlerFunc) Handle(ctx context.Context, cmd Command) ([]events.Event, error) {
	return chf(ctx, cmd)
}

// CommandHandlerMiddleware is self-described
type CommandHandlerMiddleware func(CommandHandler) CommandHandler

// CommandHandlerMultiMiddleware is self-described
func CommandHandlerMultiMiddleware(mws ...CommandHandlerMiddleware) CommandHandlerMiddleware {
	return func(ch CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]events.Event, error) {
			mw := mws[0](ch)
			for _, outerMw := range mws[1:] {
				mw = outerMw(mw)
			}
			return mw.Handle(ctx, cmd)
		})
	}
}

// ChErrMw is a command handler middleware
func ChErrMw(l Logger) CommandHandlerMiddleware {
	return func(ch CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]events.Event, error) {
			evs, err := ch.Handle(ctx, cmd)
			if err != nil {
				b, _ := json.Marshal(cmd)
				l.Printf("ch, name: %s, command: %s, error: %s\n", cmd.Name(), string(b), err.Error())
			}
			return evs, err
		})
	}
}

// Bus is an Events bus
type Bus interface {
	Dispatch(context.Context, bus.Dispatchable) (interface{}, error)
}

// ChEventMw is a domain events handler middleware
func ChEventMw(eventsBus Bus) CommandHandlerMiddleware {
	return func(ch CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]events.Event, error) {
			evs, err := ch.Handle(ctx, cmd)
			if err == nil {
				for _, e := range evs {
					eventsBus.Dispatch(ctx, e)
				}
			}
			return evs, err
		})
	}
}

// Query is a CQRS query
type Query interface {
	Name() string
}

// QueryResult is self-described
type QueryResult interface{}

// QueryHandler handles a query
type QueryHandler interface {
	Handle(ctx context.Context, q Query) (QueryResult, error)
}

// QueryHandlerFunc is a function that implements QueryHandler interface
type QueryHandlerFunc func(ctx context.Context, q Query) (QueryResult, error)

// Handle implements the QueryHandler interface
func (chf QueryHandlerFunc) Handle(ctx context.Context, q Query) (QueryResult, error) {
	return chf(ctx, q)
}

// QueryHandlerMiddleware is self-described
type QueryHandlerMiddleware func(QueryHandler) QueryHandler

// QueryHandlerMultiMiddleware is self-described
func QueryHandlerMultiMiddleware(mws ...QueryHandlerMiddleware) QueryHandlerMiddleware {
	return func(ch QueryHandler) QueryHandler {
		return QueryHandlerFunc(func(ctx context.Context, cmd Query) (QueryResult, error) {
			mw := mws[0](ch)
			for _, outerMw := range mws[1:] {
				mw = outerMw(mw)
			}
			return mw.Handle(ctx, cmd)
		})
	}
}

// QhErrMw is a query handler middleware
func QhErrMw(l Logger) QueryHandlerMiddleware {
	return func(ch QueryHandler) QueryHandler {
		return QueryHandlerFunc(func(ctx context.Context, cmd Query) (QueryResult, error) {
			evs, err := ch.Handle(ctx, cmd)
			if err != nil {
				b, _ := json.Marshal(cmd)
				l.Printf("ch, name: %s, command: %s, error: %s\n", cmd.Name(), string(b), err.Error())
			}
			return evs, err
		})
	}
}
