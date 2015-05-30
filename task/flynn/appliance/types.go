package appliance

import (
	"github.com/ory-am/gitdeploy/task"
	"time"
)

type Appliance interface {
	Attach(args ...interface{}) (id string, w *task.WorkerLog, env map[string]string, err error)
	Destroy(id string) (err error)
}

type ReleaseContainer struct {
	Manifest string
	URL      string
}

type Manifest struct {
	Env       map[string]string      `json:"env,omitempty"`
	Processes map[string]ProcessType `json:"processes"`
}

type ProcessType struct {
	Cmd         []string          `json:"cmd,omitempty"`
	Entrypoint  []string          `json:"entrypoint,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Ports       []Port            `json:"ports,omitempty"`
	Data        bool              `json:"data,omitempty"`
	Omni        bool              `json:"omni,omitempty"` // omnipresent - present on all hosts
	HostNetwork bool              `json:"host_network,omitempty"`
	Service     string            `json:"service,omitempty"`
	Resurrect   bool              `json:"resurrect,omitempty"`
}

type Port struct {
	Port    int      `json:"port"`
	Proto   string   `json:"proto"`
	Service *Service `json:"service,omitempty"`
}

type Service struct {
	Name string `json:"name,omitempty"`
	// Create the service in service discovery
	Create bool         `json:"create,omitempty"`
	Check  *HealthCheck `json:"check,omitempty"`
}

type HealthCheck struct {
	// Type is one of tcp, http, https
	Type string `json:"type,omitempty"`
	// Interval is the time to wait between checks after the service has been
	// marked as up. It defaults to two seconds.
	Interval time.Duration `json:"interval,omitempty"`
	// Threshold is the number of consecutive checks of the same status before
	// a service will be marked as up or down after coming up for the first
	// time. It defaults to 2.
	Threshold int `json:"threshold,omitempty"`
	// If KillDown is true, the job will be killed if the service goes down (or
	// does not come up)
	KillDown bool `json:"kill_down,omitempty"`
	// StartTimeout is the maximum duration that a service can take to come up
	// for the first time if KillDown is true. It defaults to ten seconds.
	StartTimeout time.Duration `json:"start_timeout,omitempty"`

	// Extra optional config fields for http/https checks
	Path   string `json:"path,omitempty"`
	Host   string `json:"host,omitempty"`
	Match  string `json:"match,omitempty"`
	Status int    `json:"status.omitempty"`
}
