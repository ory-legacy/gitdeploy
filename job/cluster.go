package job

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"net/url"
	"os/exec"
	"regexp"
)

// Clone runs the "git clone" job.
func GetCluster(em *event.EventManager, app string) (*url.URL, error) {
	reg := regexp.MustCompile(`(?mi)[a-z0-9\-A-Z]+\s+(https\:\/\/)controller\.([\.a-zA-Z0-9]+)\s+\(default\)$`)
	eventName := "jobs.cluster"
	em.Trigger(eventName, gde.New(app, "Looking up cluster..."))
	o, err := exec.Command("flynn", "cluster").CombinedOutput()
	if err != nil {
		return nil, err
	}
	if results := reg.FindStringSubmatch(string(o)); len(results) < 2 {
		em.Trigger(eventName, gde.New(app, "INTERNAL ERROR: Could not parse cluster information"))
		return nil, errors.New(fmt.Sprintf("Could not parse cluster information. Result: %s. Data: %s", results, o))
	} else {
		if u, err := url.Parse(results[1] + results[2]); err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not parse cluster information: %s", results[1])))
			return nil, err
		} else {
			return u, nil
		}
	}
}
