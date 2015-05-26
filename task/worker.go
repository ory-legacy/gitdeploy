package task

import (
	"fmt"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/event"
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
	err error
}

func (w *WorkerLog) Add(event, message string) {
	n := append(*w, &workerLogEntry{
		event:   event,
		message: message,
	})
	w = &n
}

func (w *WorkerLog) AddError(event string, err error) {
	n := append(*w, &workerLogEntry{
		event:   event,
		message: fmt.Sprintf("An error occured: %s", err.Error()),
		err: err,
	})
	w = &n
}

func (l *EventManagerWorker) Work(m []Command) error {
	for _, c := range m {
		events, err := c.Run()
		if err != nil {
			return err
		}
		for _, v := range events {
			l.EventManager.TriggerAndWait(v.event, sse.NewEvent(l.Channel, v.message))
		}
	}
	return nil
}
