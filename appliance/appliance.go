package appliance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"io/ioutil"
	"os"
	"runtime"
)

type Appliance interface {
	Attach(args ...interface{}) (id string, w *task.WorkerLog, env map[string]string, err error)
	Destroy(id string) (err error)
}

type Manifest struct {
	Processes map[string]ServerManifest `json:"processes"`
}

type PortsManifest struct {
	Port    int             `json:"port"`
	Proto   string          `json:"proto"`
	Service *ServerManifest `json:"service"`
}

type ServiceManifest struct {
	Name   string         `json:"name"`
	Create bool           `json:"create"`
	Check  *CheckManifest `json:"check"`
}

type CheckManifest struct {
	Type string `json:"type"`
}

type ServerManifest struct {
	Cmd   []string        `json:"cmd"`
	Data  bool            `json:"data"`
	Ports []PortsManifest `json:"ports"`
}

func CreateManifest(id, cmd, process string, port int, data bool) (string, error) {
	m := &Manifest{
		Processes: map[string]ServerManifest{
			process: {
				Cmd:  cmd,
				Data: data,
				Ports: []PortsManifest{
					{
						Port:  port,
						Proto: "tcp",
						Service: &ServiceManifest{
							Name:   id,
							Create: true,
							Check: CheckManifest{
								Type: "tcp",
							},
						},
					},
				},
			},
		},
	}
	manifestPath := createDirectory(id) + "manifest.json"
	if enc, err := json.MarshalIndent(m, "", "\t"); err != nil {
		return "", errors.New(fmt.Sprintf("Could not marshall manifest: %s", err.Error()))
	} else if err := ioutil.WriteFile(manifestPath, enc, 0644); err != nil {
		return "", errors.New(fmt.Sprintf("Could not write manifset: %s", err.Error()))
	}
	return manifestPath, nil
}

func createDirectory(id string) (destination string) {
	destination = fmt.Sprintf("%s/%s/", os.TempDir(), id)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s\\", os.TempDir(), id)
	}
	return destination
}
