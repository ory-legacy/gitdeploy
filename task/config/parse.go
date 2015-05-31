package config

import (
	"errors"
	"fmt"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func Parse(wd string, f func(*Config)) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
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

		f(c)
		return nil
	}
}

func ParseGodir(config *Config, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
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
		return nil
	}
}

func ParseBuildpack(config *Config, f *flynn.EnvHelper) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if len(config.Buildpack) > 0 {
			w.Add(fmt.Sprintf("Found custom buildpack url: %s.", config.Buildpack))
			f.AddEnvVar("BUILDPACK_URL", config.Buildpack)
		}
		return nil
	}
}

func ParseEnv(config *Config, f *flynn.EnvHelper) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if len(config.Env) > 0 {
			for k, v := range config.Env {
				w.Add(fmt.Sprintf("Found env var %s=%s", k, v))
				f.AddEnvVar(k, v)
			}
		}
		return nil
	}
}

func ParseProcs(config *Config, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
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
}