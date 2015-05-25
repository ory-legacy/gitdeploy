package task

import (
	"github.com/ory-am/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createEventManager() *event.EventManager {
	return event.New()
}

type listener struct {
	called map[string]interface{}
}

func (m *listener) Trigger(event string, data interface{}) {
	m.called[event] = data
}

type command struct{}

func (c *command) Run() (WorkerLog, error) {
	return WorkerLog{
		&workerLogEntry{"eventA", "abc"},
		&workerLogEntry{"eventB", "def"},
		&workerLogEntry{"eventC", "hij"},
	}, nil
}

func TestLogWorkerWork(t *testing.T) {
	em := createEventManager()
	m := make(map[string]interface{})
	l := &listener{called: m}
	em.AttachListener("eventA", l)
	em.AttachListener("eventB", l)

	lw := &EventManagerWorker{EventManager: em, Channel: "channel"}

	lw.Work([]Command{new(command)})
	assert.Equal(t, len(m), 2)
}
