package actor

import (
	"context"
	"encoding/json"
	"go.temporal.io/sdk/workflow"
	"time"
)

type Command struct {
	Name   string
	UserId string
}

type Event struct {
	Name   string
	UserId string
}

type State struct {
	Score int
}

const (
	CommandSignalName = "commands"

	TaskQueueName = "actor-sample"
)

func Workflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    2 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Actor workflow started")

	signalChannel := workflow.GetSignalChannel(ctx, CommandSignalName)
	s := workflow.NewSelector(ctx)

	var events []Event
	var state State

	// Here we expose the current event history. We could also expose "events_after" to allow consumers to pull events from an offset they track
	if err := workflow.SetQueryHandler(ctx, "current_events", func() ([]Event, error) {
		return events, nil
	}); err != nil {
		return err
	}

	// Here we expose the current state. We could expose one query per view we care about
	if err := workflow.SetQueryHandler(ctx, "current_state", func() (State, error) {
		return state, nil
	}); err != nil {
		return err
	}

	var commandString string
	s.AddReceive(signalChannel, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &commandString)
		logger.Info("Received Command", "command", commandString)

		var command Command
		json.Unmarshal([]byte(commandString), &command)
		logger.Info("Unmarshalled Command", "command", command)

		var newEvents []Event
		err := workflow.ExecuteActivity(ctx, CommandHandler, events, command).Get(ctx, &newEvents)

		// Here we ignore invalid commands
		if err != nil {
			logger.Error("Command failed.", "Error", err)
		}
		logger.Info("Processed Command", "events", newEvents)

		events = append(events, newEvents...)
		// Here we update in-memory state given new events (currentState, newEvents) => newState
		state = updateState(state, newEvents)
		// Here we might publish new events to event bus / trigger a signal to a workflow. Otherwise, they can lazily query the event stream
	})
	for {
		s.Select(ctx)
	}
}

func updateState(state State, events []Event) State {
	for _, e := range events {
		switch e.Name {
		case "upvoted":
			state.Score = state.Score + 1
			break
		case "downvoted":
			state.Score = state.Score - 1
			break
		}
	}
	return state
}

func CommandHandler(ctx context.Context, events []Event, command Command) ([]Event, error) {
	switch command.Name {
	case "upvote":
		return []Event{
			{Name: "upvoted", UserId: command.UserId},
		}, nil
		break
	case "downvote":
		return []Event{
			{Name: "downvoted", UserId: command.UserId},
		}, nil
	}
	return []Event{}, nil
}
