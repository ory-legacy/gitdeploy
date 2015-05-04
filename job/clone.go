package job

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"os"
	"os/exec"
	"runtime"
)

// Clone runs the "git clone" job.
func Clone(em *event.EventManager, app, source string) (destination string, err error) {
	eventName := "jobs.clone"
	destination = fmt.Sprintf("%s/%s", os.TempDir(), app)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s", os.TempDir(), app)
	}

	em.Trigger(eventName, gde.New(app, "Starting git clone..."))
	e := exec.Command("git", "clone", "--progress", source, destination)

	// Set up the pipes. We need to catch both stdout and stderr because git writes to both.
	stdout, err := e.StdoutPipe()
	handlePipeErr(err, em, eventName, app)
	stderr, err := e.StderrPipe()
	handlePipeErr(err, em, eventName, app)

	// Create a scanner and forward the messages to the EventManager
	go scanPipe(stdout, em, eventName, app)
	go scanPipe(stderr, em, eventName, app)

	// Run the job!
	return destination, e.Run()
}
