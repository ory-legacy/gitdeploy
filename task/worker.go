package task

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
)

type Command interface {
	Run() (WorkerLog, error)
}

type EventManagerWorker struct {
	EventManager *event.EventManager
	Channel      string
}

type WorkerLog []*workerLogEntry

type workerLogEntry struct {
	event   string
	message string
}

func (w *WorkerLog) Add(event, message string) {
	n := append(*w, &workerLogEntry{
		event:   event,
		message: message,
	})
	w = &n
}

func (w *WorkerLog) AddError(event, err error) {
	w.Add(event, fmt.Sprintf("An error occured: %s", err.Error()))
}

func (l *EventManagerWorker) Work(m []Command) error {
	for _, c := range m {
		events, err := c.Run()
		if err != nil {
			return err
		}
		for _, v := range events {
			l.EventManager.TriggerAndWait(v.event, gde.New(l.Channel, v.message))
		}
	}
	return nil
}
