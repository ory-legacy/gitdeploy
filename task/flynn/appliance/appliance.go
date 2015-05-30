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

const processName = "process"

// Create deploys a docker image.
// Workflow taken from https://gist.github.com/lmars/8be1952a8d03f8a31b17
func Create(w task.WorkerLog, id, manifestPath, url string, port int) (err error) {
	if err = flynn.CreateApp(id, "")(w); err != nil {
		return
	}

	// r := flynn.CreateReleaseContainer(manifest, "url://tbd", id, eventName, wd)
	if err = flynn.ReleaseContainer(id, manifestPath, url)(w); err != nil {
		return
	}

	if err = flynn.ScaleApp(id, processName)(w); err != nil {
		return
	}
	return
}

func CreateManifest(id string, port int, cmd []string) (string, error) {
	m := &Manifest{
		Processes: map[string]ProcessType{
			processName: {
				Cmd:  cmd,
				Data: false,
				Ports: []Port{
					{
						Port:  port,
						Proto: "tcp",
						Service: &Service{
							Name:   id,
							Create: true,
							Check: &HealthCheck{
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
