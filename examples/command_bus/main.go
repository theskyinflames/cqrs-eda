package main

import (
	"context"
	"fmt"
	"time"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"

	"github.com/google/uuid"
)

// AddUserCommand is a command
type AddUserCommand struct {
	ID       uuid.UUID
	UserName string
}

const addUserCommandName = "add_user"

// Name implements cqrs.Name interface
func (ac AddUserCommand) Name() string {
	return addUserCommandName
}

// AddUserCommandHandler is a command handler
type AddUserCommandHandler struct{}

// Handle implements cqrs.CommandHandler interface
func (ch AddUserCommandHandler) Handle(ctx context.Context, cmd cqrs.Command) ([]events.Event, error) {
	addUserCmd, ok := cmd.(AddUserCommand)
	if !ok {
		return nil, fmt.Errorf("expected command %s, but received %s", addUserCommandName, cmd.Name())
	}
	fmt.Printf("added user: %s (%s)\n", addUserCmd.UserName, addUserCmd.ID)
	return nil, nil
}

func main() {
	bus := bus.New()
	bus.Register(addUserCommandName, helpers.BusChHandler(AddUserCommandHandler{}))

	cmd := AddUserCommand{
		ID:       uuid.New(),
		UserName: "Bond, James Bond",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Dispatch(ctx, cmd)

	// Give time to output traces
	time.Sleep(time.Second)
}
