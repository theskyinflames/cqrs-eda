package cqrs

//go:generate moq -stub -out mock_cqrs_test.go -pkg cqrs_test . Command Query CommandHandler QueryHandler Event
//go:generate moq -stub -out mock_logger_test.go -pkg cqrs_test . Logger

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

// Logger is an interface
type Logger interface {
	Printf(format string, v ...interface{})
}

// Event is an event
type Event interface {
	Name() string
	AggregateID() uuid.UUID
}

// Command is a CQRS command
type Command interface {
	Name() string
}

// CommandHandler handles a command
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) ([]Event, error)
}

// CommandHandlerFunc is a function that implements CommandHandler interface
type CommandHandlerFunc func(ctx context.Context, cmd Command) ([]Event, error)

// Handle implements the CommandHandler interface
func (chf CommandHandlerFunc) Handle(ctx context.Context, cmd Command) ([]Event, error) {
	return chf(ctx, cmd)
}

// CommandHandlerMiddleware is self-described
type CommandHandlerMiddleware func(CommandHandler) CommandHandler

// ChErrMw is a command handler middleware
func ChErrMw(l Logger) CommandHandlerMiddleware {
	return func(ch CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]Event, error) {
			evs, err := ch.Handle(ctx, cmd)
			if err != nil {
				b, _ := json.Marshal(cmd)
				l.Printf("ch, name: %s, command: %s, error: %s\n", cmd.Name(), string(b), err.Error())
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
