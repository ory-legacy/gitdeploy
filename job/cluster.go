package job

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"os/exec"
	"strings"
)

// Clone runs the "git clone" job.
func GetCluster(em *event.EventManager, app string) (string, error) {
	eventName := "jobs.cluster"
	em.Trigger(eventName, gde.New(app, "Looking up cluster..."))
	o, err := exec.Command("flynn", "cluster", "default").CombinedOutput()
	if err != nil {
		return "", err
	}
	if cluster, ok := strings.Split(string(o), `"`)[1]; !ok {
		em.Trigger(eventName, gde.New(app, "INTERNAL ERROR: Could not parse cluster information"))
		return "", errors.New(fmt.Sprintf("Could not parse cluster information: %s", o))
	} else {
		return cluster, nil
	}
}
