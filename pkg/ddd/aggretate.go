package ddd

import (
	"sync"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// AggregateBasic implements
type AggregateBasic struct {
	id     uuid.UUID
	events []events.Event

	mux *sync.Mutex
}

// NewAggregateBasic is a constructor
func NewAggregateBasic(ID uuid.UUID) AggregateBasic {
	return AggregateBasic{id: ID, mux: &sync.Mutex{}}
}

// ID is a getter
func (ab AggregateBasic) ID() uuid.UUID {
	return ab.id
}

// RecordEvent is self-described
func (ab *AggregateBasic) RecordEvent(e events.Event) {
	ab.mux.Lock()
	defer ab.mux.Unlock()
	ab.events = append(ab.events, e)
}

// Events is self-described
func (ab *AggregateBasic) Events() []events.Event {
	ab.mux.Lock()
	defer ab.mux.Unlock()
	e := ab.events
	ab.events = []events.Event{}
	return e
}
