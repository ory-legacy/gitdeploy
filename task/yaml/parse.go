package yaml

import (
	"fmt"
	"errors"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/git"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func Parse(w task.WorkerLog, f *flynn.EnvHelper, app, wd string) func(task.WorkerLog) (error) {
	return func(w task.WorkerLog) error {
		w.Add("Parsing .gitdeploy.yml...")

		filename := wd + "/.gitdeploy.yml"
		if _, err := os.Stat(filename); err != nil {
			w.Add("WARN: .gitdeploy.yml not found, using defaults.")
			return nil
		}

		c := new(Config)
		if dat, err := ioutil.ReadFile(filename); err != nil {
			return errors.New(fmt.Sprintf("Could not open .gitdeploy.yml: %s", err.Error()))
		} else if err = yaml.Unmarshal(dat, c); err != nil {
			return errors.New(fmt.Sprintf("Could not parse .gitdeploy.yml: %s", err.Error()))
		}

		if err := godir(w, c, wd); err != nil {
			return w, err
		}
		if err := procs(w, c, wd); err != nil {
			return w, err
		}
		if err := buildpack(w, c, f); err != nil {
			return w, err
		}
		env(w, c, f)
		if err := f.CommitEnvVars(app); err != nil {
			return w, err
		}
		if err := git.AddAll(w, wd); err != nil {
			return w, err
		}
		if err := git.Commit(w, wd); err != nil {
			return w, err
		}

		return w, nil
	}
}

func godir(w *task.WorkerLog, config *Config, wd string) error {
	if config.Godir != "" {
		godirPath := wd + "/.godir"
		if _, err := os.Stat(godirPath); err == nil {
			w.Add("WARNING: overriding existing .godir.")
			if err := os.Remove(godirPath); err != nil {
				return errors.New(fmt.Sprintf("Could not remove existing .godir: %s", err.Error()))
			}
		}

		pm := []byte(config.Godir)
		if err := ioutil.WriteFile(godirPath, pm, 0644); err != nil {
			return errors.New(fmt.Sprintf("Could not create .godir: %s", err.Error()))
		}
	} else {
		if _, err := os.Stat(wd + "/Godeps"); err != nil {
			if _, err := os.Stat(wd + "/.godir"); err != nil {
				w.Add("WARNING: Found neither .godir config, .godir file nor Godeps directory. If this is a Go application, deployment will fail.")
			}
		}
	}
	return false, nil
}

func buildpack(w *task.WorkerLog, config *Config, f *flynn.EnvHelper) error {
	if len(config.Buildpack) > 0 {
		w.Add(fmt.Sprintf("Found custom buildpack url: %s.", config.Buildpack))
		f.AddEnvVar("BUILDPACK_URL", config.Buildpack)
	}
	return nil
}

func env(w *task.WorkerLog, config *Config, f *flynn.EnvHelper) {
	if len(config.Env) > 0 {
		for k, v := range config.Env {
			w.Add(fmt.Sprintf("Found env var %s=%s", k, v))
			f.AddEnvVar(k, v)
		}
	}
}

func procs(w *task.WorkerLog, config *Config, wd string) error {
	if len(config.ProcConfig) > 0 {
		procfilePath := wd + "/Procfile"
		if _, err := os.Stat(procfilePath); err == nil {
			w.Add("WARNING: overriding existing Procfile.")
			if err := os.Remove(procfilePath); err != nil {
				return errors.New(fmt.Sprintf("Could not remove existing Procfile: %s", err.Error()))
			}
		}

		pm, err := yaml.Marshal(config.ProcConfig)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not parse procs section: %s", err.Error()))
		}

		err = ioutil.WriteFile(procfilePath, pm, 0644)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not create Procfile: %s", err.Error()))
		}
	} else {
		if _, err := os.Stat(wd + "/Procfile"); err != nil {
			w.Add("WARNING: No procs config and no Procfile present. Deployment might fail.")
		}
	}
	return nil
}
