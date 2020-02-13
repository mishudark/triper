package basicevent

import (
	"errors"
	"github.com/golang/glog"
	"time"

	"github.com/mishudark/triper"
)


var (
	// ErrSubjectMissing When an event has no subject id
	ErrSubjectMissing = errors.New("SubjectId not found")
	// ErrCreatedAtMissing When an event has no cratedAt date
	ErrCreatedAtMissing = errors.New("CreatedAt not found")
)

// BasicEvent is the base structure for the events
type BasicEvent struct {
	triper.BaseAggregate
	SubjectId string
	Payload   map[string]string
	Source    string
	Tag       string
	Target    string
	CreatedAt time.Time
}

// Assign the values to the corresponding fields on an event
func (event *BasicEventCreated) CleanEvent() (*BasicEventCreated, error) {

	if val, ok := event.Payload["subject_id"]; ok {
		event.SubjectId = val
		delete(event.Payload, "subject_id")
	} else {
		return nil, ErrSubjectMissing
	}

	if val, ok := event.Payload["created_at"]; ok {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			glog.Fatalln(err)
		}
		event.CreatedAt = t
		delete(event.Payload, "created_at")
	} else {
		return nil, ErrCreatedAtMissing
	}

	if val, ok := event.Payload["source"]; ok {
		event.Source = val
		delete(event.Payload, "source")
	}

	if val, ok := event.Payload["target"]; ok {
		event.Target = val
		delete(event.Payload, "target")
	}

	return event, nil
}

// Reduce applies the change to the Event
func (ev *BasicEvent) Reduce(event triper.Event) error {
	switch event.Data.(type) {
	case *BasicEventCreated:
		ev.ID = event.AggregateID
	case *ProductChanged:
		ev.ID = event.AggregateID
	case *CustomerCreated:
		ev.ID = event.AggregateID
	default:
		return errors.New("undefined event")
	}
	return nil
}

// HandleCommand create events and validate based on such command
func (ev *BasicEvent) HandleCommand(command triper.Command) error {
	event := triper.Event{
		AggregateID:   ev.ID,
		AggregateType: "BasicEvent",
	}

	switch c := command.(type) {
	case *CreateCustomer:
		event.AggregateID = c.AggregateID

		baseEvent := BasicEventCreated{
			Payload: c.Payload,
			Tag:     c.Type,
		}
		newEvent, err := baseEvent.CleanEvent()
		if err != nil {
			return err
		}
		event.Data = &CustomerCreated{
			newEvent,
		}

		event.Data = newEvent

	case *ChangeProduct:
		event.AggregateID = c.AggregateID
		baseEvent := BasicEventCreated{
			Payload: c.Payload,
			Tag:     c.Type,
		}
		event.Data = &ProductChanged{
			&baseEvent,
		}

	}

	err := ev.Reduce(event)
	if err != nil {
		return err
	}
	return nil
}
