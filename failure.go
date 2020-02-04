package triper

import "fmt"

// FailureType defines the alert(error) type while a command is being processed
type FailureType string

// nolint
const (
	FailureLoadingEvents     FailureType = "loading_events"
	FailureReplayingEvents   FailureType = "replying_events"
	FailureProcessingCommand FailureType = "processing_command"
	FailureInvalidID         FailureType = "invalid_id"
	FailureSavingOnStorage   FailureType = "saving_on_storage"
	FailurePublishingEvents  FailureType = "publishing_events"
	FailureVersionMissmatch  FailureType = "version_missmatch"
)

// Failure is an error while the command is being processed
type Failure struct {
	CommandID      string      `json:"command_id"`
	CommandType    string      `json:"command_type"`
	CommandVersion int         `json:"command_version"`
	AggregateID    string      `json:"aggregate_id"`
	AggregateType  string      `json:"aggregate_type"`
	Type           FailureType `json:"type"`
	Err            error       `json:"error"`
}

// NewFailure returns an alert that implements an error interface
func NewFailure(err error, typ FailureType, command Command) error {
	if err == nil {
		return nil
	}

	return Failure{
		CommandID:      command.GetID(),
		CommandType:    command.GetType(),
		CommandVersion: command.GetVersion(),
		AggregateID:    command.GetAggregateID(),
		AggregateType:  command.GetAggregateType(),
		Type:           typ,
		Err:            err,
	}
}

func (f Failure) Error() string {
	return fmt.Sprintf("[%s]: command-id=%s command-version=%d aggregate-id=%s error=%s",
		f.Type,
		f.CommandID,
		f.CommandVersion,
		f.AggregateID,
		f.Err)
}
