package helpers

import (
	"context"
	"errors"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

func BusChHandler(ch cqrs.CommandHandler) bus.Handler {
	return func(ctx context.Context, d bus.Dispatchable) (interface{}, error) {
		cmd, ok := d.(cqrs.Command)
		if !ok {
			return nil, errors.New("unexpected dispatchable")
		}
		return ch.Handle(ctx, cmd)
	}
}

func BusQhHandler(ch cqrs.QueryHandler) bus.Handler {
	return func(ctx context.Context, d bus.Dispatchable) (interface{}, error) {
		cmd, ok := d.(cqrs.Query)
		if !ok {
			return nil, errors.New("unexpected dispatchable")
		}
		return ch.Handle(ctx, cmd)
	}
}
