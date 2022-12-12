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

/*
This is an example of concurrent events bus use case.
*/
func main() {
	listeners, inChs := listeners(2)
	busHandlers := busHandlers(inChs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the events listeners
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			fmt.Println(err.Error())
		}
	}()
	for _, l := range listeners {
		go l.Listen(ctx, errChan)
	}

	// Create the events bus
	const (
		concurrencyLimit   = 2
		dispatchingTimeout = time.Second
	)
	bus := bus.NewConcurrentBus(dispatchingTimeout, concurrencyLimit)

	// Register the handlers for each event to be dispatched to their listeners
	for i, bh := range busHandlers {
		bus.Register(eventName(i), bh)
	}
	go bus.Run(ctx)

	// Dispatch an event 1
	rsChan := bus.Dispatch(ctx, events.NewEventBasic(uuid.New(), eventName(0), nil))

	// Dispatch an event 2
	rsChan2 := bus.Dispatch(ctx, events.NewEventBasic(uuid.New(), eventName(1), nil))

	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case eventHndResponse := <-rsChan:
				fmt.Printf("RS 1: %#v\n", eventHndResponse)
				i++
			case eventHndResponse := <-rsChan2:
				fmt.Printf("RS 2: %#v\n", eventHndResponse)
				i++
			}
			if i == len(inChs) {
				return
			}
		}
	}()

	// Giving time to output traces
	time.Sleep(time.Second)
}

func listeners(n int) ([]events.Listener, []chan events.Event) {
	var (
		listeners []events.Listener
		inChan    []chan events.Event
	)
	for i := 0; i < n; i++ {
		i := i
		eventsChan := make(chan events.Event)
		eventHandlers := []events.Handler{
			func(e events.Event) {
				fmt.Printf("eh %d.1, received %s event, with id %s\n", i, e.Name(), e.AggregateID().String())
			},
			func(e events.Event) {
				fmt.Printf("eh %d.2, received %s event, with id %s\n", i, e.Name(), e.AggregateID().String())
			},
		}
		listeners = append(listeners, events.NewListener(eventsChan, eventName(i), eventHandlers...))
		inChan = append(inChan, eventsChan)
	}

	return listeners, inChan
}

func eventName(i int) string {
	return fmt.Sprintf("event%d", i)
}

func busHandlers(inChs []chan events.Event) []bus.Handler {
	var busHnds []bus.Handler
	for i := range inChs {
		i := i
		busHnds = append(busHnds, func(_ context.Context, d bus.Dispatchable) (interface{}, error) {
			e, ok := d.(events.Event)
			if !ok {
				return nil, errors.New("unexpected dispatchable")
			}
			inChs[i] <- e
			return fmt.Sprintf("dispatchable processed by hnd %d\n", i), nil
		})
	}
	return busHnds
}
