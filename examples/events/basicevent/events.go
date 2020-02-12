package basicevent

import "time"

//Thi is the base Event Created
type BasicEventCreated struct {
	SubjectId string            `json:"subject_id"`
	Payload   map[string]string `json:"payload"`
	Source    string            `json:"source"`
	Target    string            `json:"target"`
	CreatedAt time.Time         `json:"created_at"`
	Tag       string            `json:"tag"`
}

type CustomerCreated struct {
	Event *BasicEventCreated `json:"event"`
}

type ProductChanged struct {
	Event *BasicEventCreated `json:"event"`
}
