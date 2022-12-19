package ddd

import (
	"sync"

	"github.com/google/uuid"
)

// Event is self-described
type Event interface {
	Name() string
}

type Aggregate interface {
	ID() uuid.UUID
	RecordEvent(e Event)
	Events() []Event
}

type AggregateBasic struct {
	sync.Mutex

	ID     uuid.UUID
	events []Event
}

func (ab *AggregateBasic) RecordEvent(e Event) {
	ab.Lock()
	defer ab.Unlock()
	ab.events = append(ab.events, e)
}

func (ab *AggregateBasic) Events() []Event {
	ab.Lock()
	defer ab.Unlock()
	e := ab.events
	ab.events = []Event{}
	return e
}
