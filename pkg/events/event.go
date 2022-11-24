package events

import "github.com/google/uuid"

// EventBasic is a domain event
type EventBasic struct {
	ID          uuid.UUID
	aggregateID uuid.UUID
	name        string
	body        any
}

// NewEventBasic is a constructor
func NewEventBasic(aggregateID uuid.UUID, name string, body any) EventBasic {
	return EventBasic{
		ID:          uuid.New(),
		aggregateID: aggregateID,
		name:        name,
		body:        body,
	}
}

// Name is a getter
func (e EventBasic) Name() string {
	return e.name
}

// AggregateID is a getter
func (e EventBasic) AggregateID() uuid.UUID {
	return e.aggregateID
}

// Body is a getter
func (e EventBasic) Body() any {
	return e.body
}
