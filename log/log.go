package log

import (
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"log"
)

type Listener struct{}

func (l *Listener) Trigger(event string, data interface{}) {
	if e, ok := data.(gde.JobEvent); ok {
		// TODO Ugly...
		e.SetEventName(event)
		log.Printf("Log listener: Event %s on app %s said: %s", event, e.GetApp(), e.GetMessage())
		return
	}
	log.Fatalf("Log listener: Type mismatch: %s is not job.JobEvent", data)
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
