package job

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
)

type DeployConfig struct {
	Procfile ProcfileConfig `yaml:"Procfile"`
	Godir    string         `yaml:"Godir"`
}

type ProcfileConfig map[string]string

// FIXME
var eventName = "jobs.parse"

func Parse(em *event.EventManager, app, sourcePath string) error {
	filename := sourcePath + "/.gitdeploy.yml"
	if _, err := os.Stat(filename); err != nil {
		em.Trigger(eventName, gde.New(app, "WARN: .gitdeploy.yml not found. Deployment might to fail."))
		return nil
	}

	em.Trigger(eventName, gde.New(app, "Found .gitdeploy.yml"))
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		em.Trigger(eventName, gde.New(app, fmt.Sprintf("ERROR: Could not open .gitdeploy.yml: %s", err.Error())))
		return err
	}
	config := new(DeployConfig)
	if err = yaml.Unmarshal(dat, config); err != nil {
		em.Trigger(eventName, gde.New(app, fmt.Sprintf("ERROR: Could not parse .gitdeploy.yml: %s", err.Error())))
		return err
	}
	em.Trigger(eventName, gde.New(app, "Successfully parsed .gitdeploy.yml"))

	// Generate Procfile
	commit := false
	if created, err := genProcfile(em, config, app, sourcePath); err != nil {
		return err
	} else if created {
		commit = true
	}
	if created, err := genGodir(em, config, app, sourcePath); err != nil {
		return err
	} else if created {
		commit = true
	}

	if commit {
		em.Trigger(eventName, gde.New(app, "Commiting changes..."))
        em.Trigger(eventName, gde.New(app, "Setting git user..."))
        if err := runGitCommand(em, app, sourcePath, "config", "user.name", "gitdeploy"); err != nil {
            em.Trigger(eventName, gde.New(app, "ERROR in git config user.name!"))
            return err
        }
        if err := runGitCommand(em, app, sourcePath, "config", "user.email", "hello@gitdeploy.io"); err != nil {
            em.Trigger(eventName, gde.New(app, "ERROR in git config user.email!"))
            return err
        }
		if err := runGitCommand(em, app, sourcePath, "commit", "-a", "-m", "gitdeploy"); err != nil {
			em.Trigger(eventName, gde.New(app, "ERROR in git commit!"))
			return err
		}
	}

	return nil
}

func genGodir(em *event.EventManager, config *DeployConfig, app, sourcePath string) (bool, error) {
	em.Trigger(eventName, gde.New(app, fmt.Sprintf("%s", config.Godir)))
	if config.Godir != "" {
		godirPath := sourcePath + "/.godir"
		em.Trigger(eventName, gde.New(app, "Found .godir configuration..."))

		// Remove existing Procfile
		if _, err := os.Stat(godirPath); err == nil {
			em.Trigger(eventName, gde.New(app, "Overriding existing .godir..."))
			if err := os.Remove(godirPath); err != nil {
				em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not remove existing .godir: %s", err.Error())))
				return false, err
			}
		}

		pm := []byte(config.Godir)
		if err := ioutil.WriteFile(godirPath, pm, 0644); err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not create .godir: %s", err.Error())))
			return false, err
		}

		em.Trigger(eventName, gde.New(app, ".godir created!"))
		if err := runGitCommand(em, app, sourcePath, "add", ".godir"); err != nil {
			em.Trigger(eventName, gde.New(app, "ERROR in git add!"))
			return false, err
		}
		return true, nil
	} else {
		em.Trigger(eventName, gde.New(app, "WARNING: No .godir configuration found. If this is a Go app deployment might fail."))
	}
	return false, nil
}

func genProcfile(em *event.EventManager, config *DeployConfig, app, sourcePath string) (bool, error) {
	if len(config.Procfile) > 0 {
		procfilePath := sourcePath + "/Procfile"
		em.Trigger(eventName, gde.New(app, "Found Procfile configuration..."))

		// Remove existing Procfile
		if _, err := os.Stat(procfilePath); err == nil {
			em.Trigger(eventName, gde.New(app, "Overriding existing Procfile..."))
			if err := os.Remove(procfilePath); err != nil {
				em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not remove existing Procfile: %s", err.Error())))
				return false, err
			}
		}

		pm, err := yaml.Marshal(config.Procfile)
		if err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not marshall configuration: %s", err.Error())))
			return false, err
		}

		err = ioutil.WriteFile(procfilePath, pm, 0644)
		if err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not create Procfile: %s", err.Error())))
			return false, err
		}

		em.Trigger(eventName, gde.New(app, "Procfile created!"))
		em.Trigger(eventName, gde.New(app, fmt.Sprintf("%s", config.Procfile)))
		if err := runGitCommand(em, app, sourcePath, "add", "Procfile"); err != nil {
			em.Trigger(eventName, gde.New(app, "ERROR in git add!"))
			return false, err
		}
		return true, nil
	} else {
		em.Trigger(eventName, gde.New(app, "WARNING: No Procfile configuration found. Deployment is most likely going to fail."))
	}
	return false, nil
}

func runGitCommand(em *event.EventManager, app, sourcePath string, args ...string) error {
	e := exec.Command("git", args...)
	e.Dir = sourcePath

	// Set up the pipes. We need to catch both stdout and stderr because git writes to both.
	stdout, err := e.StdoutPipe()
	handlePipeErr(err, em, eventName, app)
	stderr, err := e.StderrPipe()
	handlePipeErr(err, em, eventName, app)

	// Create a scanner and forward the messages to the EventManager
	go scanPipe(stdout, em, eventName, app)
	go scanPipe(stderr, em, eventName, app)

	return e.Run()
}
