package task

import (
	"fmt"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"time"
)

// Task is a function which runs a task, for example "flynn create app" or "git clone"
type Task func(w WorkerLog) (error)

// WorkerLog
type WorkerLog chan *workerEvent

type workerEvent struct {
	message string
	err     error
	// offset should be int32 or float64
	offset  time.Time
}

// Add adds an event to the channel
func (w WorkerLog) Add(message string) {
	w <- &workerEvent{
		message: message,
		err:     nil,
		offset: time.Now(),
	}
}

// Add adds an error event to the channel
func (w WorkerLog) AddError(err error) {
	w <- &workerEvent{
		message: fmt.Sprintf("An error occured: %s", err.Error()),
		err:     err,
		offset: time.Now(),
	}
}

// RunJob ...
func RunJob(event, channel string, em *event.EventManager, taskList []Task) (err error) {
	// Run all jobs
	for _, task := range taskList {
		workerChan := make(WorkerLog)
		errChan := make(chan error)

		// Run next in line
		go func() {
			defer close(workerChan)
			defer close(errChan)
			errChan <- task(workerChan)
		}()

		// Forward output to the event manager
		go func() {
			for {
				select {
				case v, open := <-workerChan:
					if !open {
						return
					}
					em.TriggerAndWait(event, sse.NewEvent(channel, v.message, event))
				}
			}
		}()

		// Catch errors if any occur
		for {
			select {
			case err := <-errChan:
				return err
			}
		}
	}
	return nil
}
