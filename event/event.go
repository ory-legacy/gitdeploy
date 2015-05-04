package event

import (
	"encoding/json"
)

type JobEvent interface {
	GetMessage() string
	GetApp() string
	SetEventName(eventName string)
}

type jobEventLog struct {
	Message   string `json:"data"`
	App       string `json:"app"`
	EventName string `json:"eventName"`
}

func New(app, message string) *jobEventLog {
	j := new(jobEventLog)
	j.Message = message
	j.App = app
	return j
}

func (j *jobEventLog) GetMessage() string {
	r, _ := json.Marshal(j)
	return string(r)
}

func (j *jobEventLog) GetApp() string {
	return j.App
}

func (j *jobEventLog) SetEventName(eventName string) {
	j.EventName = eventName
}
