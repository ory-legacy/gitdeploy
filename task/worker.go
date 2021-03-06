package task

import (
	"fmt"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"log"
	"time"
)

// Task is a function which runs a task, for example "flynn create app" or "git clone"
type Task func(w WorkerLog) error

// TaskList is a list of Tasks
type TaskList struct {
	Tasks []TaskManifest
}

type TaskManifest struct {
	Name string
	Task Task
}

// WorkerLog
type WorkerLog chan *workerEvent

type workerEvent struct {
	message string
	err     error
	// offset should be int32 or float64
	offset time.Time
}

func (tl *TaskList) Add(event string, task Task) {
	if tl.Tasks == nil {
		tl.Tasks = make([]TaskManifest, 0)
	}
	tl.Tasks = append(tl.Tasks, TaskManifest{
		Task: task,
		Name: event,
	})
}

// Add adds an event to the channel
func (w WorkerLog) Add(message string) {
	go func(w WorkerLog, message string) {
		fmt.Printf("trying to write: %s", message)
		w <- &workerEvent{
			message: message,
			err:     nil,
			offset:  time.Now(),
		}
	}(w, message)
}

// Add adds an error event to the channel
func (w WorkerLog) AddError(err error) {
	go func(w WorkerLog, err error) {
		w <- &workerEvent{
			message: fmt.Sprintf("An error occured: %s", err.Error()),
			err:     err,
			offset:  time.Now(),
		}
	}(w, err)
}

// RunJob sequentially runs all tasks in the TaskList using the FIFO scheduling algorithm.
// Output generated by a task are forwarded to the event manager and eventually
// to all listeners, like the storage backend or the logger.
// When one of the tasks returns an error, RunJob will abort all queued tasks and return the error.
func RunJob(channel string, em *event.EventManager, taskList *TaskList) error {
	var c WorkerLog
	cs := make([]WorkerLog, 0)
	defer func() {
		go func() {
			// Output from a task is read in a goroutine. If the process exists,
			// the goroutine could try to write something, although we've already closed
			// the channel which results in a panic.
			// To circumvent this, we've added a short sleep which should make sure, that all
			// output has been read.
			time.Sleep(60 * time.Second)
			for _, v := range cs {
				close(v)
			}
		}()
	}()
	for _, task := range taskList.Tasks {
		c = make(WorkerLog)
		cs = append(cs, c)
		go scanTask(em, c, task, channel)
		err := task.Task(c)
		if err != nil {
			em.Trigger("error", sse.NewEvent(channel, err.Error(), "error"))
			return err
		}
	}
	return nil
}

func scanTask(em *event.EventManager, c WorkerLog, t TaskManifest, channel string) {
	defer log.Printf("Exiting event manager listener for %s", t.Name)
	ev := t.Name
	for {
		select {
		case v, open := <-c:
			if !open {
				return
			}
			em.Trigger(ev, sse.NewEvent(channel, v.message, ev))
		}
	}
}
