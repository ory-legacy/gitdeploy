package flynn

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
    gexec "github.com/ory-am/gitdeploy/exec"
	"os/exec"
)

const deployEventName = "jobs.deploy"

type KeyAdd struct {
    *Helper
}

func (d *KeyAdd) Run() (gexec.WorkerLog, error) {
    w := new(gexec.WorkerLog)
    w.Add("jobs.deploy", "Adding key...")

    if err := d.Helper.exec(w, "flynn", "key", "add"); err != nil {
        return w, err
    }
    return w, nil
}

func Deploy() error {
	deployEventName := "jobs.deploy"
	em.Trigger(deployEventName, gde.New(app, "Starting with deployment..."))

	// Add flynn key
	if err := runFlynnKeyAdd(em, deployEventName, app, sourcePath); err != nil {
		return err
	}
	// Create app container on flynn
	if err := runFlynnCreateApp(em, deployEventName, app, sourcePath); err != nil {
		return err
	}
	// Push app to flynn
	if err := runFlynnPush(em, deployEventName, app, sourcePath); err != nil {
		return err
	}

	return nil
}

func (d *Deploy) runFlynnKeyAdd(w *gexec.WorkerLog) error {
	return nil
}

func runFlynnCreateApp(em *event.EventManager, deployEventName string, app string, sourcePath string) error {
	em.Trigger(deployEventName, gde.New(app, "Creating app..."))
    w.Add("jobs.deploy", "Creating app...")
	if err := run(em, deployEventName, app, sourcePath, "flynn", "create", "-y", app); err != nil {
		return err
	}
	return nil
}

func runFlynnPush(em *event.EventManager, deployEventName string, app string, sourcePath string) error {
	em.Trigger(deployEventName, gde.New(app, "Pushing app..."))
	if err := run(em, deployEventName, app, sourcePath, "git", "push", "flynn", "master", "--progress"); err != nil {
		return err
	}
	return nil
}

func run(em *event.EventManager, deployEventName, app, sourcePath string, cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Dir = sourcePath

	// Set up the pipes. We need to catch both stdout and stderr because git writes to both.
	// Errors can be rejected because handlePipeErr fatals on errors.
	stdout, err := e.StdoutPipe()
	handlePipeErr(err, em, deployEventName, app)
	stderr, err := e.StderrPipe()
	handlePipeErr(err, em, deployEventName, app)

	// Create a scanner and forward the messages to the EventManager
	go scanPipe(stdout, em, deployEventName, app)
	go scanPipe(stderr, em, deployEventName, app)

	err = e.Run()
	if err != nil {
		em.Trigger(deployEventName, gde.New(app, fmt.Sprintf("An error occured: %s", err.Error())))
		return err
	}

	return err
}
