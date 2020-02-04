package main

import (
	"EventLogTest/eventlog"
	"flag"
	"github.com/golang/glog"
	"github.com/mishudark/triper"
	"log"
	"os"
	"time"
)

func main() {
	flag.Parse()

	commandBus, err := GetConfig()
	if err != nil {
		glog.Infoln(err)
		os.Exit(1)
	}

	end := make(chan bool)

	for i := 0; i < 1; i++ {
		go func() {
			uuid := triper.GenerateUUID()

			//1) Create Event on log
			newEvent := eventlog.CreateLogevent{
					Owner : "valerybriz",
					Payload : map[string]string{
						"apple":      "pen",
						"source":     "Finance",
						"pinaple":    "ppap",
						"subject_id": "123232",
						"created_at": "2020-01-30T15:04:05Z",
						"target":     "Accountability",
					},
			}
			newEvent.AggregateID = uuid
			newEvent.Type = "create-logevent"

			commandBus.HandleCommand(&newEvent)
			log.Println("event created", uuid)

			//2) Change event
			time.Sleep(time.Millisecond * 100)
			eventChange := eventlog.ChangeLogevent{
				Owner : "valerybriz",
				Payload : map[string]string{
					"apple":      "pen",
					"source":     "CRM",
					"pinaple":    "ppap",
					"subject_id": "123232",
					"created_at": "2020-03-30T19:04:05Z",
					"target":     "CRM",
				},
			}
			eventChange.AggregateID = uuid
			eventChange.Type = "change-on-event"

			commandBus.HandleCommand(&eventChange)
			log.Println("event changed", uuid)


		}()
	}
	<-end
}
