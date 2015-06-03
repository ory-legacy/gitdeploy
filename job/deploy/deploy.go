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
        return createAppliances(conf, f, func(name, id string) error {
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
        w.Add(app.ID)
        return nil
    })
    return tasks
}

func createAppliances(c *config.Config, eh *flynn.EnvHelper, f func(name, id string) error) func(w task.WorkerLog) error {
    return func(w task.WorkerLog) error {
        for name, conf := range c.Addons {
            id := uuid.NewRandom().String()
            switch strings.ToLower(name) {
                case "mongodb":
                    w.Add(fmt.Sprintf("Attaching mongodb with id %s", id))
                    mongo.Create(id, &conf, eh)(w)
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
