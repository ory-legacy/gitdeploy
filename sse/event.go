package sse
import "encoding/json"

// An SSE event
type Event struct {
	Message   string `json:"data"`
	App       string `json:"app"`
	EventName string `json:"eventName"`
}

func NewEvent(app, message string) *Event {
	j := new(Event)
	j.Message = message
	j.App = app
	return j
}

func (j *Event) SSEify() string {
	r, _ := json.MarshalIndent(j, "data: ", "\t")
	return string(r)
}

func (j *Event) GetApp() string {
	return j.App
}

func (j *Event) SetEventName(eventName string) {
	j.EventName = eventName
}