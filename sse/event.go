package sse

import (
	"encoding/json"
	"fmt"
	"log"
)

// An SSE event
type Event struct {
	Data      string `json:"data"`
	App       string `json:"app"`
	EventName string `json:"eventName"`
}

func NewEvent(app, data, eventName string) *Event {
	j := new(Event)
	j.Data = data
	j.App = app
	j.EventName = eventName
	return j
}

func (j *Event) SSEify() string {
	r, err := json.Marshal(j)
	if err != nil {
		msg := fmt.Sprintf("Could not marshall %s: %s", j, err.Error())
		log.Println(msg)
		return "data: " + msg
	}
	return "data: " + string(r)
}
