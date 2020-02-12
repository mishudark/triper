package main

import (
	"EventLogTest/basicevent"
	"errors"
	"flag"
	"github.com/golang/glog"
	"github.com/mishudark/triper"
	"log"
	"os"
	"time"
)

func Categorize(commandBus triper.CommandBus, event basicevent.BasicEvent) error {
	switch event.Tag {
	case "customer-created":
		customerEvent := basicevent.CreateCustomer{
			Payload: event.Payload,
		}
		customerEvent.Type = event.Tag
		customerEvent.AggregateID = event.ID
		commandBus.HandleCommand(&customerEvent)

		return nil
	case "product-changed":
		changeProductEvent := basicevent.ChangeProduct{
			Payload: event.Payload,
		}
		changeProductEvent.Type = event.Tag
		changeProductEvent.AggregateID = event.ID
		commandBus.HandleCommand(&changeProductEvent)

		return nil
	default:
		return errors.New("undefined event")
	}
}

func main() {
	flag.Parse()

	commandBus, err := GetConfig()
	if err != nil {
		glog.Infoln(err)
		os.Exit(1)
	}

	end := make(chan bool)

	//1) Product change
	productChange := basicevent.BasicEvent{
		// The following map simulates a JSON payload coming from an external source
		Payload: map[string]string{
			"source":     "Finance",
			"subject_id": "123232",
			"created_at": "2020-01-30T15:04:05Z",
			"target":     "Accountability",
			"approved":   "yes",
			"tag":        "product-changed",
		},
	}

	//2) Customer created
	time.Sleep(time.Millisecond * 100)
	customerCreated := basicevent.BasicEvent{
		// The following map simulates a JSON payload coming from an external source
		Payload: map[string]string{
			"customer_name": "Miguel",
			"source":        "CRM",
			"date_of_birth": "1980-03-30T19:04:05Z",
			"subject_id":    "5467646",
			"created_at":    "2020-03-30T19:04:05Z",
			"target":        "CRM",
			"tag":           "customer-created",
		},
	}

	events := []basicevent.BasicEvent{customerCreated, productChange}

	uuid := triper.GenerateUUID()

	for _, currentEvent := range events {

		if eventTag, ok := currentEvent.Payload["tag"]; ok {

			go func() { // here starts the go routine
				delete(currentEvent.Payload, "tag")
				currentEvent.Tag = eventTag
				currentEvent.ID = uuid
				err := Categorize(commandBus, currentEvent)
				if err != nil {
					log.Fatalf("Error Categorizing the event %s", uuid)
				}
				glog.Infof("event created", uuid)
			}() // here ends the go routine

		} else {
			log.Fatalf("Couldn't find a valid TAG for the event %s", uuid)
		}

	}
	<-end
}
