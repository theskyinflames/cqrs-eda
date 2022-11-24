package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/events"

	"github.com/google/uuid"
)

const eventName = "anEvent"

/*
This is an example of events bus use case.
*/
func main() {
	var (
		eventsChan    = make(chan events.Event)
		eventHandlers = []events.Handler{
			func(e events.Event) {
				fmt.Printf("eh1, received %s event, with id %s\n", e.Name(), e.AggregateID().String())
			},
			func(e events.Event) {
				fmt.Printf("eh2, received %s event, with id %s\n", e.Name(), e.AggregateID().String())
			},
		}
		eventsListener             = events.NewListener(eventsChan, eventName, eventHandlers...)
		busHandler     bus.Handler = func(_ context.Context, d bus.Dispatchable) (any, error) {
			e, ok := d.(events.Event)
			if !ok {
				return nil, errors.New("unexpected dispatchable")
			}
			eventsChan <- e
			return "dispatchable processed", nil
		}
	)

	// Start the events listener
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			fmt.Printf("events listener: %s", err.Error())
		}
	}()
	go eventsListener.Listen(ctx, errChan)

	// Start the events bus
	bus := bus.New()
	bus.Register(eventName, busHandler)

	// Dispatch an event
	bus.Dispatch(ctx, events.NewEventBasic(uuid.New(), eventName, nil))

	// Give time to output the logs
	time.Sleep(time.Second)
}
