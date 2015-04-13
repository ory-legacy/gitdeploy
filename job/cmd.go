package job

import (
    "bufio"
    "fmt"
    "github.com/ory-am/event"
    "io"
    "log"
    gde "github.com/ory-am/gitdeploy/event"
)

// Handle pipe errors. Exits fatally when an error occurred.
func handlePipeErr(err error, em *event.EventManager, eventName, app string) {
    if err != nil {
        m := fmt.Sprintf("Error creating a pipe. This should not happen. Message was: %s", err)
        em.Trigger(eventName, gde.New(app, m))
        log.Fatalf(m)
    }
}

// Scan the pipe and forward the messages to the EventManager
func scanPipe(p io.ReadCloser, em *event.EventManager, eventName, app string) {
    s := bufio.NewScanner(p)
    for s.Scan() {
        em.Trigger(eventName, gde.New(app, s.Text()))
    }
}
