package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type eventListener struct{}

const eventName = "foobar"

type eventData struct {
	t *testing.T
	c chan bool
	m string
}

func (l *eventListener) Trigger(event string, data interface{}) {
	if d, ok := data.(*eventData); ok {
		d.c <- true
		return
	}
	panic("Type assertion failed")
}

type dontListen struct{}

func (l *dontListen) Trigger(event string, data interface{}) {
	if c, ok := data.(*eventData); ok {
		c.t.FailNow()
		return
	}
	panic("Type assertion failed")
}

type waitListener struct{}

func (l *waitListener) Trigger(event string, data interface{}) {
	if c, ok := data.(*eventData); ok {
		c.m = "foobar"
		return
	}
	panic("Type assertion failed")
}

func TestTriggerAndWait(t *testing.T) {
	data := &eventData{t, nil, ""}
	em := New()
	el := new(waitListener)

	em.AttachListener(eventName, el)
	em.TriggerAndWait(eventName, data)

	assert.Equal(t, "foobar", data.m)
}

func TestTrigger(t *testing.T) {
	data := &eventData{t, make(chan bool), ""}
	em := New()
	el := new(eventListener)

	em.AttachListener(eventName, el)
	em.Trigger(eventName, data)

	assert.True(t, <-data.c)
}

func TestDetach(t *testing.T) {
	data := &eventData{t, make(chan bool), ""}
	em := New()
	el := new(dontListen)

	em.AttachListener(eventName, el)
	em.DetachListener(eventName, el)
	em.Trigger(eventName, data)

	assert.True(t, true)
}
