package yaml
import (
	"github.com/ory-am/gitdeploy/task"
	"os"
	"io/ioutil"
	"github.com/ory-am/gitdeploy/Godeps/_workspace/src/github.com/go-errors/errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task/git"
)

type Parse struct{
	*task.Helper
	Flynn *flynn.Flynn
	Git   *git.Git
}

func (d *Parse) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(d.EventName, "Parsing .gitdeploy.yml...")

	filename := d.WorkingDirectory + "/.gitdeploy.yml"
	if _, err := os.Stat(filename); err != nil {
		w.Add(d.EventName, "WARN: .gitdeploy.yml not found, using defaults.")
		return nil
	}

	c := new(Config)
	if dat, err := ioutil.ReadFile(filename); err != nil {
		return errors.New(fmt.Sprintf("Could not open .gitdeploy.yml: %s", err.Error()))
	} else if err = yaml.Unmarshal(dat, c); err != nil {
		return errors.New(fmt.Sprintf("Could not parse .gitdeploy.yml: %s", err.Error()))
	}

	if err := d.godir(w, c); err != nil {
		return w, err
	}
	if err := d.procs(w, c); err != nil {
		return w, err
	}
	if err := d.buildpack(w, c); err != nil {
		return w, err
	}
	d.env(w, c)
	if err := d.Flynn.CommitEnvVars(); err != nil {
		return w, err
	}
	if err := d.Git.AddAll(); err != nil {
		return w, err
	}
	if err := d.Git.Commit(); err != nil {
		return w, err
	}

	return w, nil
}

func (h *Parse) godir(w *task.WorkerLog, config *Config) (error) {
	if config.Godir != "" {
		godirPath := h.WorkingDirectory + "/.godir"
		if _, err := os.Stat(godirPath); err == nil {
			w.Add(h.EventName, "WARNING: overriding existing .godir.")
			if err := os.Remove(godirPath); err != nil {
				return errors.New(fmt.Sprintf("Could not remove existing .godir: %s", err.Error()))
			}
		}

		pm := []byte(config.Godir)
		if err := ioutil.WriteFile(godirPath, pm, 0644); err != nil {
			return errors.New(fmt.Sprintf("Could not create .godir: %s", err.Error()))
		}
	} else {
		if _, err := os.Stat(h.WorkingDirectory + "/Godeps"); err != nil {
			if _, err := os.Stat(h.WorkingDirectory + "/.godir"); err != nil {
				w.Add(h.EventName, "WARNING: Found neither .godir config, .godir file nor Godeps directory. If this is a Go application, deployment will fail.")
			}
		}
	}
	return false, nil
}

func (h *Parse) buildpack(w *task.WorkerLog, config *Config) (error) {
	if len(config.Buildpack) > 0 {
		w.Add(h.EventName, fmt.Sprintf("Found custom buildpack url: %s.", config.Buildpack))
		h.Flynn.AddEnvVar("BUILDPACK_URL", config.Buildpack)
	}
	return nil
}

func (h *Parse) env(w *task.WorkerLog, config *Config) {
	if len(config.Env) > 0 {
		for k, v := range config.Env {
			w.Add(h.EventName, fmt.Sprintf("Found env var %s=%s", k, v))
			h.Flynn.AddEnvVar(k, v)
		}
	}
}

func (h *Parse) procs(w *task.WorkerLog, config *Config) (error) {
	if len(config.ProcConfig) > 0 {
		procfilePath := h.WorkingDirectory + "/Procfile"
		if _, err := os.Stat(procfilePath); err == nil {
			w.Add(h.EventName, "WARNING: overriding existing Procfile.")
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
		if _, err := os.Stat(h.WorkingDirectory + "/Procfile"); err != nil {
			w.Add(h.EventName, "WARNING: No procs config and no Procfile present. Deployment might fail.")
		}
	}
	return nil
}