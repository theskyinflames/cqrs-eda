package bus

//go:generate moq -stub -out mock_bus_test.go -pkg bus_test . Dispatchable

import (
	"context"
	"errors"
)

// Dispatchable is a dispatchable item through the bus
type Dispatchable interface {
	Name() string
}

// Handler handles a dispatchable from the bus
type Handler func(ctx context.Context, d Dispatchable) (interface{}, error)

// Bus is self-described
type Bus struct {
	h map[string]Handler
}

// New is a constructor
func New() Bus {
	return Bus{
		h: make(map[string]Handler),
	}
}

// Register adds a new handler to the bus
func (b Bus) Register(n string, h Handler) {
	b.h[n] = h
}

// ErrNotDispatchable is self-described
var ErrNotDispatchable = errors.New("not dispatchable")

// Dispatch dispatches a dispatchable item
func (b Bus) Dispatch(ctx context.Context, d Dispatchable) (interface{}, error) {
	h, ok := b.h[d.Name()]
	if !ok {
		return nil, ErrNotDispatchable
	}
	return h(ctx, d)
}
