package ddd_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/ddd"
)

type TestEvent struct {
	name string
}

func (e TestEvent) Name() string {
	return e.name
}

func TestRecordEventConcurrent(t *testing.T) {
	a := ddd.AggregateBasic{
		ID: uuid.New(),
	}

	// Record events concurrently
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			a.RecordEvent(TestEvent{name: "event1"})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			a.RecordEvent(TestEvent{name: "event2"})
		}
	}()
	wg.Wait()

	// Check that all events were recorded correctly
	events := a.Events()
	require.Len(t, events, 2000)

	event1Count := 0
	event2Count := 0
	for _, e := range events {
		switch e.Name() {
		case "event1":
			event1Count++
		case "event2":
			event2Count++
		default:
			t.Error("Unexpected event name:", e.Name())
		}
	}
	if event1Count != 1000 {
		t.Error("Expected 1000 event1 events, got", event1Count)
	}
	if event2Count != 1000 {
		t.Error("Expected 1000 event2 events, got", event2Count)
	}
}
