package job

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"os/exec"
)

func Deploy(em *event.EventManager, app string, sourcePath string) error {
	eventName := "jobs.deploy"
	em.Trigger(eventName, gde.New(app, "Starting with deployment..."))

	// Add flynn key
	if err := runFlynnKeyAdd(em, eventName, app, sourcePath); err != nil {
		return err
	}
	// Create app container on flynn
	if err := runFlynnCreateApp(em, eventName, app, sourcePath); err != nil {
		return err
	}
	// Push app to flynn
	if err := runFlynnPush(em, eventName, app, sourcePath); err != nil {
		return err
	}

	return nil
}

func runFlynnKeyAdd(em *event.EventManager, eventName string, app string, sourcePath string) error {
	em.Trigger(eventName, gde.New(app, "Adding key..."))
	if err := run(em, eventName, app, sourcePath, "flynn", "key", "add"); err != nil {
		return err
	}
	return nil
}

func runFlynnCreateApp(em *event.EventManager, eventName string, app string, sourcePath string) error {
	em.Trigger(eventName, gde.New(app, "Creating app..."))
	if err := run(em, eventName, app, sourcePath, "flynn", "create", "-y", app); err != nil {
		return err
	}
	return nil
}

func runFlynnPush(em *event.EventManager, eventName string, app string, sourcePath string) error {
	em.Trigger(eventName, gde.New(app, "Pushing app..."))
	if err := run(em, eventName, app, sourcePath, "git", "push", "flynn", "master", "--progress"); err != nil {
		return err
	}
	return nil
}

func run(em *event.EventManager, eventName, app, sourcePath string, cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Dir = sourcePath

	// Set up the pipes. We need to catch both stdout and stderr because git writes to both.
	// Errors can be rejected because handlePipeErr fatals on errors.
	stdout, err := e.StdoutPipe()
	handlePipeErr(err, em, eventName, app)
	stderr, err := e.StderrPipe()
	handlePipeErr(err, em, eventName, app)

	// Create a scanner and forward the messages to the EventManager
	go scanPipe(stdout, em, eventName, app)
	go scanPipe(stderr, em, eventName, app)

	err = e.Run()
	if err != nil {
		em.Trigger(eventName, gde.New(app, fmt.Sprintf("An error occured: %s", err.Error())))
		return err
	}

	return err
}
