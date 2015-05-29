package mongo

import (
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn/appliance"
	"github.com/ory-am/gitdeploy/task/flynn"
)

func Create(w task.WorkerLog, id string, config map[string]string) func(w task.WorkerLog) (err error) {
	return func(w task.WorkerLog) error {
		// tbd
		url := "https://registry.hub.docker.com?name=mongo&id=216d9a0e82646f77b31b78eeb0e26db5500930bbd6085d5d5c3844ec27c0ca50"
		manifest, err := appliance.CreateManifest(id, 27017, []string{"mongod"})
		if err != nil {
			return
		}

		env, err := appliance.Create(w, id, manifest, url, 27017, config)
		if err != nil {
			return
		}

		eh := new(flynn.EnvHelper)
		for k, v := range env {
			eh.AddEnvVar(k, v)
		}
		if err = eh.CommitEnvVars(id); err != nil {
			return
		}

		return
	}
}
