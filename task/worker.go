package task

import (
	"fmt"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"time"
)

type Command interface {
	Run(w WorkerLog, e WorkerError)
}

type EventManagerWorker struct {
	EventManager *event.EventManager
	Channel      string
}

type WorkerLog chan *workerLogEntry
type WorkerError chan error

type workerLogEntry struct {
	event   string
	message string
	err     error
	timestamp time.Time
}

func (w *WorkerLog) Add(event, message string) {
	w <- &workerLogEntry{
		event:   event,
		message: message,
		timestamp: time.Now(),
	}
}

func (w *WorkerLog) AddError(event string, err error) {
	w <- &workerLogEntry{
		event:   event,
		message: fmt.Sprintf("An error occured: %s", err.Error()),
		err:     err,
		timestamp: time.Now(),
	}
}

func (l *EventManagerWorker) Work(m []Command) (err error) {
	w := make(WorkerLog, 10)
	e := make(WorkerError)

	// Run all jobs
	for _, c := range m {
		// Run next in line
		go func() {
			defer close(w)
			defer close(e)
			c.Run(w, e)
		}()

		// Forward output to the event manager
		go func() {
			for {
				select {
				case v, open := <-w:
					if !open {
						return
					}
					l.EventManager.TriggerAndWait(v.event, sse.NewEvent(l.Channel, v.message))
				}
			}
		}()

		// Catch errors if any occur
		for {
			select {
			case err := <-e:
				return err
			}
		}
	}
	return nil
}
