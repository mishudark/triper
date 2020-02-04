package bank

import "github.com/mishudark/triper"

//CreateAccount assigned to an owner
type CreateAccount struct {
	triper.BaseCommand
	Owner string
}

//PerformDeposit to a given account
type PerformDeposit struct {
	triper.BaseCommand
	Amount int
}

//ChangeOwner of an account
type ChangeOwner struct {
	triper.BaseCommand
	Owner string
}

//PerformWithdrawal to a given account
type PerformWithdrawal struct {
	triper.BaseCommand
	Amount int
}
