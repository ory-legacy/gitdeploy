package task

import (
	"github.com/ory-am/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
	"fmt"
)

var errMock = errors.New("error")

type listener struct {
	called map[string]int
}

func (m *listener) Trigger(event string, _ interface{}) {
	m.called[event]++
	fmt.Println("asdf" + event)
}

type command struct{}

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
	tl := make(TaskList)
	tl.AddTask(mockTask)
	tl.AddTask(mockTask)
	tl.AddTask(mockTask)

	em.AttachListener("a", l)
	em.AttachListener("b", l)

	err := RunJob("channel", em, tl)
	fmt.Printf("%s", l)

	err := RunJob("channel", em, tl)
	fmt.Printf("%s", l)

	assert.Equal(t, 2, len(l.called))
	assert.Equal(t, 3, l.called["a"])
	assert.Equal(t, 3, l.called["b"])
	assert.Nil(t, err)
}

func TestAddTask(t *testing.T) {
	tl := make(TaskList)
	tl.AddTask("a", mockTask)
	tl.AddTask("a", mockTask)
	tl.AddTask("b", mockTask)

	assert.Equal(t, 2, len(tl["a"]))
	assert.Equal(t, 1, len(tl["b"]))
}

func TestRunJobError(t *testing.T) {
	em := event.New()
	tl := make(TaskList)
	tl.AddTask("a", mockErrorTask)

	err := RunJob("channel", em, tl)
	assert.Equal(t, err, errMock)
}