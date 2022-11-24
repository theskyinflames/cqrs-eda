package events

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

//go:generate moq -stub -out mock_event_test.go -pkg events_test . Event

// Event is an event
type Event interface {
	Name() string
	AggregateID() uuid.UUID
}

// Handler handles an event
type Handler func(Event)

// Listener reacts to events
type Listener struct {
	ch       <-chan Event
	name     string
	handlers []Handler
}

// NewListener is a constructor
func NewListener(ch <-chan Event, name string, handlers ...Handler) Listener {
	return Listener{
		ch:       ch,
		name:     name,
		handlers: handlers,
	}
}

var ErrUnexpectedEvent = errors.New("unexpected event")

// Listen makes the listener start receiving events
func (l Listener) Listen(ctx context.Context, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-l.ch:
			if e.Name() != l.name {
				log.Printf("expected %s, received %s", l.name, e.Name())
				errCh <- ErrUnexpectedEvent
				continue // ignores not expected events
			}
			if err := l.handle(e); err != nil {
				errCh <- fmt.Errorf("listener for %s: %w", l.name, err)
			}
		}
	}
}

// Handle handles the event
func (l Listener) handle(e Event) error {
	// TODO: Think about if it's worth to parallelize it with goroutines
	for _, h := range l.handlers {
		h(e)
	}
	return nil
}
