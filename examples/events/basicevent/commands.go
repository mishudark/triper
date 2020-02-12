package basicevent

import (
	"github.com/mishudark/triper"
)

//Create the base event
type CreateBasicEvent struct {
	triper.BaseCommand
	Payload map[string]string
}

//Create the change event
type ChangeBasicEvent struct {
	triper.BaseCommand
	Payload map[string]string
}

//Create a customer
type CreateCustomer struct {
	triper.BaseCommand
	Payload map[string]string
}

//Change the information of a product
type ChangeProduct struct {
	triper.BaseCommand
	Payload map[string]string
}
