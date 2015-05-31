package deploy

import (
	"github.com/ory-am/gitdeploy/storage"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/config"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/git"
	"github.com/ory-am/gitdeploy/task/janitor"
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"fmt"
	"github.com/ory-am/gitdeploy/task/flynn/appliance/mongo"
	"errors"
)

func CreateJob(w task.WorkerLog, store storage.Storage, app *storage.App) (tasks task.TaskList) {
	var conf *config.Config
	dir := git.CreateDirectory(app.ID)
	f := new(flynn.EnvHelper)

	return task.TaskList{
		"git.clone": git.Clone(app.Repository, dir),
		"git.checkout": git.Checkout(app.ID, dir, app.Ref),
		"config.parse": config.Parse(dir, func(c *config.Config) {
			conf = c
		}),
		"config.procs":     config.ParseProcs(conf, dir),
		"config.buildpack": config.ParseBuildpack(conf, f),
		"config.env":       config.ParseEnv(conf, f),
		// TODO Refactor so we don't need task.WorkerLog here...
		// TODO Tasks should be able to add tasks to the queue.
		// TODO It should be enough for now to add the tasks as next in line
		"config.appliances": createAppliances(w, conf, f, func(name, id string) error {
			_, err := store.AddAppliance(app.ID, id, name)
			return err
		}),
		"env.commit": func(w task.WorkerLog) error {
			w.Add("Commiting env vars...")
			return f.CommitEnvVars(app.ID)
		},
		"git.add":     git.AddAll(dir),
		"git.commit":  git.Commit(dir),
		"app.release": flynn.ReleaseApp(dir),
		"app.cleanup": janitor.Cleanup(dir),
		"app.deployed": func(w task.WorkerLog) error {
			w.Add(app.ID)
			return nil
		},
	}
}

func createAppliances(w task.WorkerLog, c *config.Config, eh *flynn.EnvHelper, f func(name, id string) error) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		for name, conf := range c.Addons {
			id := uuid.NewRandom().String()
			switch strings.ToLower(name) {
			case "mongodb":
				w.Add("Attaching mongodb...")
				mongo.Create(id, &conf, eh)(w)
				return f(name, id)
			default:
				w.AddError(errors.New(fmt.Sprintf("Appliance %s not supported", name)))
			}
		}
		return nil
	}
}