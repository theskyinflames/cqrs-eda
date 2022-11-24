package events_test

import (
	"context"
	"testing"

	"github.com/theskyinflames/cqrs-eda/pkg/events"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestListenerListen(t *testing.T) {
	t.Run(`Given a not expected event, when it's listened, then it's ignored`, func(t *testing.T) {
		var (
			handlerCalls int
			handler      = func(_ events.Event) {
				handlerCalls++
			}
			name = "anEvent"
			ch   = make(chan events.Event)
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		l := events.NewListener(ch, name, handler)
		errChan := make(chan error)
		go l.Listen(ctx, errChan)

		ch <- events.NewEventBasic(uuid.New(), "unexpected", nil)

		err := <-errChan
		require.ErrorIs(t, events.ErrUnexpectedEvent, err)
		require.Equal(t, 0, handlerCalls)
	})

	t.Run(`Given an expected event, when it's listened, then it's handled`, func(t *testing.T) {
		var (
			handlerCalled = make(chan struct{})
			handler       = func(_ events.Event) {
				close(handlerCalled)
			}
			name = "anEvent"
			ch   = make(chan events.Event)
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		l := events.NewListener(ch, name, handler)
		errChan := make(chan error)

		go l.Listen(ctx, errChan)

		ch <- events.NewEventBasic(uuid.New(), "anEvent", nil)

		<-handlerCalled
	})
}
