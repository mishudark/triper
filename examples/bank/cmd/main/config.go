package main

import (
	"github.com/mishudark/triper"
	"github.com/mishudark/triper/commandhandler/basic"
	"github.com/mishudark/triper/config"
	"github.com/mishudark/triper/examples/bank"
)

func getConfig() (triper.CommandBus, error) {
	//register events
	reg := triper.NewEventRegister()
	reg.Set(bank.AccountCreated{})
	reg.Set(bank.DepositPerformed{})
	reg.Set(bank.WithdrawalPerformed{})

	//eventbus
	// rabbit, err := config.RabbitMq("guest", "guest", "localhost", 5672)

	return config.NewClient(
		config.Badger("/tmp", reg),                                  // event store
		config.Nats("nats://ruser:T0pS3cr3t@localhost:4222", false), // event bus
		config.AsyncCommandBus(30),                                  // command bus
		config.WireCommands(
			&bank.Account{},          // aggregate
			basic.NewCommandHandler,  // command handler
			"bank",                   // event store bucket
			"account",                // event store subset
			bank.CreateAccount{},     // command
			bank.PerformDeposit{},    // command
			bank.PerformWithdrawal{}, // command
		),
	)
}
