package basic

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/glog"

	"github.com/mishudark/triper"
)

// ErrInvalidID missing initial event.
var ErrInvalidID = errors.New("invalid ID, initial event missign")

// Handler contains the info to manage commands
type Handler struct {
	repository     *triper.Repository
	aggregate      reflect.Type
	bucket, subset string
}

// NewCommandHandler return a handler.
func NewCommandHandler(repository *triper.Repository, aggregate triper.AggregateHandler, bucket, subset string) triper.CommandHandler {
	return &Handler{
		repository: repository,
		aggregate:  reflect.TypeOf(aggregate).Elem(),
		bucket:     bucket,
		subset:     subset,
	}
}

// Handle a command, if any error is produced, it will be published to the errors bucket.
func (h *Handler) Handle(command triper.Command) (err error) {
	version := command.GetVersion()
	aggregate := reflect.New(h.aggregate).Interface().(triper.AggregateHandler)

	defer func() {
		if err != nil {
			glog.Errorln(err)
			er := h.repository.PublishError(err, command, h.bucket, "errors")

			if er != nil {
				glog.Errorln(er)
			}
		}
	}()

	if version != 0 {
		if err = h.repository.Load(aggregate, command.GetAggregateID()); err != nil {
			return triper.NewFailure(err, triper.FailureLoadingEvents, command)
		}

		if version != aggregate.GetVersion() {
			return triper.NewFailure(fmt.Errorf("got: %d, expected: %d", aggregate.GetVersion(), version), triper.FailureVersionMissmatch, command)
		}
	}

	// the aggregate can have errors trying to replay the previous events
	if aggregate.HasError() {
		return triper.NewFailure(aggregate.GetError(), triper.FailureReplayingEvents, command)
	}

	if err = aggregate.HandleCommand(command); err != nil {
		return triper.NewFailure(err, triper.FailureProcessingCommand, command)
	}

	// After to handle the command, the aggregate can have errors applying the new events
	if aggregate.HasError() {
		return triper.NewFailure(aggregate.GetError(), triper.FailureReplayingEvents, command)
	}

	// if not contain a valid ID,  the initial event (some like createAggreagate event) is missing
	if aggregate.GetID() == "" {
		return triper.NewFailure(ErrInvalidID, triper.FailureInvalidID, command)
	}

	// add the command id for traceability
	aggregate.AttachCommandID(command.GetID())

	// save the changes using the repository
	if err = h.repository.Save(aggregate, version); err != nil {
		return triper.NewFailure(err, triper.FailureSavingOnStorage, command)
	}

	err = h.repository.PublishEvents(aggregate, h.bucket, h.subset)
	return triper.NewFailure(err, triper.FailurePublishingEvents, command)
}
