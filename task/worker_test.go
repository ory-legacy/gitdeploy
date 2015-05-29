package task

import (
	"github.com/ory-am/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)

var errMock = errors.New("error")

type listener struct {
	called map[string]int
}

func (m *listener) Trigger(event string, _ interface{}) {
	m.called[event]++
}

type command struct {}

func mockTask(w WorkerLog) error {
	w.Add("abc")
	w.Add("def")
	w.Add("hij")
	return nil
}

func mockErrorTask(w WorkerLog) error {
	w.AddError(errMock)
	return errMock
}

func TestRunJob(t *testing.T) {
	em := event.New()
	l := &listener{called: make(map[string]int)}
	tl := []Task{mockTask, mockTask, mockTask}

	em.AttachListener("a", l)
	em.AttachListener("b", l)

	err := RunJob("a", "channel", em, tl)
	err = RunJob("b", "channel", em, tl)

	assert.Equal(t, 2, len(l.called))
	assert.Equal(t, 3, l.called["a"])
	assert.Equal(t, 3, l.called["b"])
	assert.Nil(t, err)
}

func TestRunJobError(t *testing.T) {
	em := event.New()
	tl := []Task{mockErrorTask}

	err := RunJob("a", "channel", em, tl)
	assert.Equal(t, err, errMock)
}