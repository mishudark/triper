package eventlog

import "time"

//Thi is the base Event Created
type LogeventCreated struct {
	Owner     string            `json:"owner"`
	SubjectId string            `json:"subject_id"`
	Payload   map[string]string `json:"payload"`
	Source    string            `json:"source"`
	Target    string            `json:"target"`
	CreatedAt time.Time         `json:"created_at"`
	Tag       string 			`json:"tag"`
}

type CustomerCreated struct {
	LogeventCreated
}

type EventChanged struct {
	LogeventCreated
}

//OwnerChanged event
type OwnerChanged struct {
	Owner string `json:"owner"`
}
