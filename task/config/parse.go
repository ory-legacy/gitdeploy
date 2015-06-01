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

		if c.Version == "" {
			c.Version = "0.1"
		} else if c.Version != "0.1" {
			return errors.New(fmt.Sprintf("Gitdeploy.yml version %s not supported.", c.Version))
		}

		f(c)
		return nil
	}
}

func ParseGodir(c *Config, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if c.Godir != "" {
			godirPath := wd + "/.godir"
			if _, err := os.Stat(godirPath); err == nil {
				w.Add("WARNING: overriding existing .godir.")
				if err := os.Remove(godirPath); err != nil {
					return errors.New(fmt.Sprintf("Could not remove existing .godir: %s", err.Error()))
				}
			}

			pm := []byte(c.Godir)
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

func ParseBuildpack(c *Config, f *flynn.EnvHelper) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if len(c.Buildpack) > 0 {
			w.Add(fmt.Sprintf("Found custom buildpack url: %s.", c.Buildpack))
			f.AddEnvVar("BUILDPACK_URL", c.Buildpack)
		}
		return nil
	}
}

func ParseEnv(c *Config, f *flynn.EnvHelper) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if len(c.Env) > 0 {
			for k, v := range c.Env {
				w.Add(fmt.Sprintf("Found env var %s=%s", k, v))
				f.AddEnvVar(k, v)
			}
		}
		return nil
	}
}

func ParseProcs(c *Config, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		if len(c.ProcConfig) > 0 {
			procfilePath := wd + "/Procfile"
			if _, err := os.Stat(procfilePath); err == nil {
				w.Add("WARNING: overriding existing Procfile.")
				if err := os.Remove(procfilePath); err != nil {
					return errors.New(fmt.Sprintf("Could not remove existing Procfile: %s", err.Error()))
				}
			}

			pm, err := yaml.Marshal(c.ProcConfig)
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
