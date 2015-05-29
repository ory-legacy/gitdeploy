package log

import (
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"log"
)

type Listener struct{}

func (l *Listener) Trigger(event string, data interface{}) {
	if e, ok := data.(sse.Event); ok {
		log.Printf("Log listener: Event %s on app %s said: %s", event, e.App, e.Data)
		return
	}
	log.Fatalf("Log listener: Type mismatch: %s is not job.Event", data)
}

func (l *Listener) AttachAggregate(em *event.EventManager) {
	em.AttachListener("app.created", l)
	em.AttachListener("jobs.clone", l)
	em.AttachListener("jobs.parse", l)
	em.AttachListener("jobs.deploy", l)
	em.AttachListener("jobs.cluster", l)
	em.AttachListener("app.deployed", l)
	em.AttachListener("app.cleanup", l)
}
