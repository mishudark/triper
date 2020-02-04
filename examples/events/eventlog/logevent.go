package eventlog

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mishudark/triper"
)

//ErrSubjectMissing When an event has no subject id
var ErrSubjectMissing = errors.New("SubjectId not found")
var ErrCreatedAtMissing = errors.New("CreatedAt not found")

// Base Structure for the events
type Logevent struct {
	triper.BaseAggregate
	Owner     string
	SubjectId string
	Payload   map[string]string
	Source    string
	Tag       string
	Target    string
	CreatedAt time.Time
}

//Assign the values to the corresponding fields on an event
func (event *LogeventCreated) CleanEvent() (*LogeventCreated, error) {

	if val, ok := event.Payload["subject_id"]; ok {
		event.SubjectId = val
		delete(event.Payload, "subject_id")
	} else {
		return nil, ErrSubjectMissing
	}

	if val, ok := event.Payload["created_at"]; ok {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			fmt.Println(err)
		}
		event.CreatedAt = t
		delete(event.Payload, "created_at")
	} else {
		return nil, ErrCreatedAtMissing
	}

	if val, ok := event.Payload["source"]; ok {
		event.Source = val
		event.Tag = fmt.Sprintf("%s-%s", val, event.CreatedAt)
		delete(event.Payload, "source")
	}

	if val, ok := event.Payload["target"]; ok {
		event.Target = val
		delete(event.Payload, "target")
	}

	return event, nil
}

//ApplyChange to account
func (ev *Logevent) Reduce(event triper.Event) error{
	switch e := event.Data.(type) {
	case *LogeventCreated:
		ev.Owner = e.Owner
		ev.ID = event.AggregateID
	case *EventChanged:
		ev.Owner = e.Owner
		ev.ID = event.AggregateID
	default:
		return errors.New("undefined event")
	}
	return nil
}

//HandleCommand create events and validate based on such command
func (ev *Logevent) HandleCommand(command triper.Command) error {
	event := triper.Event{
		AggregateID:   ev.ID,
		AggregateType: "Logevent",
	}

	switch c := command.(type) {
	case *CreateLogevent:
		event.AggregateID = c.AggregateID

		newEvent := &LogeventCreated{
			Owner:   c.Owner,
			Payload: c.Payload,
		}
		newEvent, err := newEvent.CleanEvent()
		if err != nil {
			return err
		}
		event.Data = newEvent
	case *ChangeLogevent:
		event.AggregateID = c.AggregateID

		event.Data = &EventChanged{
			LogeventCreated{
				Owner:   c.Owner,
				Payload: c.Payload,
			},
		}

	}

	log.Printf("created: %s", event.Data)

	err := ev.Reduce(event)
	if err != nil{
		return err
	}
	return nil
}
