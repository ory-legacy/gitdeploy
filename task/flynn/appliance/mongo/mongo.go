package mongo

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/config"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/flynn/appliance"
	"fmt"
)

func Create(id string, c *config.DatabaseConfig, f *flynn.EnvHelper) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if c.Version == "" {
			c.Version = "latest"
		}

		w.Add("Generating manifest for MongoDB...")
		url := fmt.Sprintf("https://registry.hub.docker.com?name=mongo&tag=%s", c.Version)
		manifest, err := appliance.CreateManifest(id, 27017, []string{"sh", "-c", "mkdir -p /data/db && mongod"})
		if err != nil {
			return err
		}

		w.Add("Releasing MongoDB...")
		if err := appliance.Create(w, id, manifest, url, 27017); err != nil {
			return err
		}

		w.Add("Setting up database and environment variables for MongoDB...")
		db := uuid.NewRandom().String()
		f.AddEnvVar(c.Host, id+".discoverd")
		f.AddEnvVar(c.Port, "27017")
		f.AddEnvVar(c.Database, db)
		f.AddEnvVar(c.URL, fmt.Sprintf("mongodb://%s.discoverd:27017/%s", id, db))
		f.AddEnvVar(c.User, "")
		f.AddEnvVar(c.Password, "")
		return nil
	}
}
