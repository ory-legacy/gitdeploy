package log

import (
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"log"
)

type Listener struct{}

func (l *Listener) Trigger(event string, data interface{}) {
	if e, ok := data.(*sse.Event); ok {
		log.Printf("Log listener: Event %s on app %s said: %s", event, e.App, e.Data)
		return
	}
	log.Fatalf("Log listener: Type mismatch: %s is not sse.Event", data)
}

func (l *Listener) AttachAggregate(em *event.EventManager) {
	em.AttachListener("git.clone", l)
	em.AttachListener("config.parse", l)
	em.AttachListener("config.procs", l)
	em.AttachListener("config.buildpack", l)
	em.AttachListener("config.env", l)
	em.AttachListener("config.appliances", l)
	em.AttachListener("env.commit", l)
	em.AttachListener("git.add", l)
	em.AttachListener("git.commit", l)
	em.AttachListener("app.release", l)
	em.AttachListener("app.create", l)
	em.AttachListener("app.cleanup", l)
	em.AttachListener("error", l)
	em.AttachListener("app.deployed", l)
}
