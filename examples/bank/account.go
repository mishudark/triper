package bank

import (
	"errors"

	"github.com/mishudark/triper"
)

// ErrBalanceOut when you don't have balance to perform the operation
var ErrBalanceOut = errors.New("balance out")

// Account of bank
type Account struct {
	triper.BaseAggregate
	Owner   string
	Balance int
}

// ApplyChange to account
func (a *Account) Reduce(event triper.Event) error {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.Owner = e.Owner
		a.ID = event.AggregateID
	case *DepositPerformed:
		a.Balance += e.Amount
	case *WithdrawalPerformed:
		a.Balance -= e.Amount
	default:
		return errors.New("undefined event")
	}

	return nil
}

// HandleCommand create events and validate based on such command
func (a *Account) HandleCommand(command triper.Command) error {
	event := triper.Event{
		AggregateID:   a.ID,
		AggregateType: "Account",
	}

	switch c := command.(type) {
	case *CreateAccount:
		event.AggregateID = c.AggregateID
		event.Data = &AccountCreated{c.Owner}

	case *PerformDeposit:
		event.Data = &DepositPerformed{
			c.Amount,
		}

	case *PerformWithdrawal:
		if a.Balance < c.Amount {
			return ErrBalanceOut
		}

		event.Data = &WithdrawalPerformed{
			c.Amount,
		}
	}

	triper.ReduceHelper(a, event, true)
	return nil
}
