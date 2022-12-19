package ddd

import (
	"sync"

	"github.com/google/uuid"
)

// Event is self-described
type Event interface {
	Name() string
}

// AggregateBasic implements
type AggregateBasic struct {
	sync.Mutex

	id     uuid.UUID
	events []Event
}

// NewAggregateBasic is a constructor
func NewAggregateBasic(ID uuid.UUID) AggregateBasic {
	return AggregateBasic{id: ID}
}

// ID is a getter
func (ab *AggregateBasic) ID() uuid.UUID {
	return ab.id
}

// RecordEvent is self-described
func (ab *AggregateBasic) RecordEvent(e Event) {
	ab.Lock()
	defer ab.Unlock()
	ab.events = append(ab.events, e)
}

// Events is self-described
func (ab *AggregateBasic) Events() []Event {
	ab.Lock()
	defer ab.Unlock()
	e := ab.events
	ab.events = []Event{}
	return e
}
