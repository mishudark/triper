package eventlog

import (
	"github.com/mishudark/triper"
)

//Create the base event
type CreateLogevent struct {
	triper.BaseCommand
	Owner   string
	Payload map[string]string
}

//Create the base event
type ChangeLogevent struct {
	triper.BaseCommand
	Owner   string
	Payload map[string]string
}

//ChangeOwner of an account
type ChangeOwner struct {
	triper.BaseCommand
	Owner string
}
