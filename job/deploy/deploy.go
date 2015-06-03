package deploy

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/ory-am/gitdeploy/storage"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/config"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/flynn/appliance/mongo"
	"github.com/ory-am/gitdeploy/task/git"
	//"github.com/ory-am/gitdeploy/task/janitor"
	"strings"
	"github.com/ory-am/gitdeploy/eco"
)

func CreateJob(store storage.Storage, app *storage.App) (tasks *task.TaskList) {
	var conf *config.Config
	dir := git.CreateDirectory(app.ID)
	f := new(flynn.EnvHelper)

	tasks = new(task.TaskList)
	tasks.Add("git.clone", git.Clone(app.Repository, dir))
	tasks.Add("git.checkout", git.Checkout(app.ID, dir, app.Ref))
	tasks.Add("app.create", flynn.CreateApp(app.ID, dir, false))
	tasks.Add("config.parse", config.Parse(dir, func(c *config.Config) {
		conf = c
	}))
	tasks.Add("config.procs", func(w task.WorkerLog) error {
		fmt.Printf("%s", conf)
		return config.ParseProcs(conf, dir)(w)
	})
	tasks.Add("config.buildpack", func(w task.WorkerLog) error {
		return config.ParseBuildpack(conf, f)(w)
	})
	tasks.Add("config.env", func(w task.WorkerLog) error {
		return config.ParseEnv(conf, f)(w)
	})
	tasks.Add("config.appliances", func(w task.WorkerLog) error {
		w.Add("Looking for appliances...")
		return createAppliances(conf, f, func(name, id string) error {
			w.Add(fmt.Sprintf("Created appliance %s: %s", name, app.ID))
			_, err := store.AddAppliance(app.ID, id, name)
			return err
		})(w)
	})
	tasks.Add("env.commit", func(w task.WorkerLog) error {
		w.Add("Commiting env vars...")
		return f.CommitEnvVars(app.ID)
	})
	tasks.Add("git.add", git.AddAll(dir))
	tasks.Add("git.commit", git.Commit(dir))
	tasks.Add("app.release", flynn.ReleaseApp(dir, app.ID))
	// tasks.Add("app.cleanup", janitor.Cleanup(dir))
	tasks.Add("app.deployed", func(w task.WorkerLog) error {
		host, err := eco.GetFlynnHost()
		if err != nil {
			return err
		}
		// Send before database query and avoid panic when
		// db is not responsive
		w.Add(app.ID)
		app.URL = fmt.Sprintf("%s.%s", app.ID, host)
		if err := store.UpdateApp(app); err != nil {
			return err
		}
		return nil
	})
	return tasks
}

func createAppliances(c *config.Config, eh *flynn.EnvHelper, f func(name, id string) error) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		for _, conf := range c.Addons {
			name := conf.Type
			w.Add(fmt.Sprintf("Iterating appliance %s", name))
			id := uuid.NewRandom().String()
			switch strings.ToLower(name) {
			case "mongodb":
				w.Add(fmt.Sprintf("Creating mongodb: %s", id))
				err := mongo.Create(id, &conf, eh)(w)
				if err != nil {
					return err
				}
				return f(name, id)
			default:
				err := errors.New(fmt.Sprintf("Appliance %s not supported", name))
				w.AddError(err)
				return err
			}
		}
		return nil
	}
}
