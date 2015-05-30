package mongo

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/config"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/flynn/appliance"
)

func Create(id string, c *config.DatabaseConfig, f *flynn.EnvHelper) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		url := "https://registry.hub.docker.com?name=mongo&id=216d9a0e82646f77b31b78eeb0e26db5500930bbd6085d5d5c3844ec27c0ca50"
		manifest, err := appliance.CreateManifest(id, 27017, []string{"mongod"})
		if err != nil {
			return err
		}

		if err := appliance.Create(w, id, manifest, url, 27017); err != nil {
			return err
		}

		db := uuid.NewRandom().String()
		f.AddEnvVar(c.Host, id+".discoverd")
		f.AddEnvVar(c.Port, "27017")
		f.AddEnvVar(c.Database, db)
		f.AddEnvVar(c.URL, "mongodb://"+id+":27017/"+db)
		f.AddEnvVar(c.User, "")
		f.AddEnvVar(c.Password, "")

		if err := f.CommitEnvVars(id); err != nil {
			return err
		}
		return nil
	}
}
