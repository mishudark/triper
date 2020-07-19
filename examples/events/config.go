package main

import (
	"EventLogTest/basicevent"
	"github.com/mishudark/triper"
	"github.com/mishudark/triper/commandhandler/basic"
	"github.com/mishudark/triper/config"
)

func GetConfig() (triper.CommandBus, error) {
	//register events
	reg := triper.NewEventRegister()
	reg.Set(basicevent.ProductChanged{})
	reg.Set(basicevent.CustomerCreated{})
	reg.Set(basicevent.BasicEventCreated{})

	//event bus
	// rabbit, err := config.RabbitMq("guest", "guest", "localhost", 5672)

	return config.NewClient(
		config.Badger("/tmp", reg), // event store
		config.Nats("nats://ruser:T0pS3cr3t@localhost:4222", false), // event bus
		config.AsyncCommandBus(30),                                  // command bus
		config.WireCommands(
			&basicevent.BasicEvent{},      // aggregate
			basic.NewCommandHandler,     // command handler
			"eventstore",                // event store bucket
			"events",                    // event store subset
			basicevent.CreateBasicEvent{}, // command
			basicevent.ChangeProduct{},
			basicevent.CreateCustomer{},
		),
	)
}
