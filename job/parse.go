package job

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type DeployConfig struct {
	Procfile ProcfileConfig `yaml:"Procfile"`
	Godir    string         `yaml:"Godir"`
}

type ProcfileConfig map[string]string

var eventName = "jobs.parse"

func Parse(em *event.EventManager, app, sourcePath string) error {
	filename := sourcePath + "/.gitdeploy.yml"
	if _, err := os.Stat(filename); err != nil {
		em.Trigger(eventName, gde.New(app, "ERROR: .gitdeploy.yml not found. Deployment is probably going to fail."))
		return err
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

	log.Printf("%s %s", dat, config)
	em.Trigger(eventName, gde.New(app, fmt.Sprintf("%s", config)))
	em.Trigger(eventName, gde.New(app, fmt.Sprintf("Godir: %s", config.Godir)))
	em.Trigger(eventName, gde.New(app, fmt.Sprintf("Procfile: %s", config.Procfile)))

	// Generate Procfile
	if err := genProcfile(em, config, app, sourcePath); err != nil {
		return err
	}
	if err := genGodir(em, config, app, sourcePath); err != nil {
		return err
	}

	return nil
}

func genGodir(em *event.EventManager, config *DeployConfig, app, sourcePath string) error {
	em.Trigger(eventName, gde.New(app, fmt.Sprintf("%s", config.Godir)))
	if config.Godir != "" {
		godirPath := sourcePath + "/.godir"
		em.Trigger(eventName, gde.New(app, "Found .godir configuration..."))

		// Remove existing Procfile
		if _, err := os.Stat(godirPath); err == nil {
			em.Trigger(eventName, gde.New(app, "Overriding existing .godir..."))
			if err := os.Remove(godirPath); err != nil {
				em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not remove existing .godir: %s", err.Error())))
				return err
			}
		}

		pm := []byte(config.Godir)
		if err := ioutil.WriteFile(godirPath, pm, 0644); err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not create .godir: %s", err.Error())))
			return err
		}

		em.Trigger(eventName, gde.New(app, ".godir created!"))
	} else {
		em.Trigger(eventName, gde.New(app, "WARNING: No .godir configuration found. If this is a Go app deployment might fail."))
	}
	return nil
}

func genProcfile(em *event.EventManager, config *DeployConfig, app, sourcePath string) error {
	if len(config.Procfile) > 0 {
		procfilePath := sourcePath + "/Procfile"
		em.Trigger(eventName, gde.New(app, "Found Procfile configuration..."))

		// Remove existing Procfile
		if _, err := os.Stat(procfilePath); err == nil {
			em.Trigger(eventName, gde.New(app, "Overriding existing Procfile..."))
			if err := os.Remove(procfilePath); err != nil {
				em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not remove existing Procfile: %s", err.Error())))
				return err
			}
		}

		pm, err := yaml.Marshal(config.Procfile)
		if err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not marshall configuration: %s", err.Error())))
			return err
		}

		err = ioutil.WriteFile(procfilePath, pm, 0644)
		if err != nil {
			em.Trigger(eventName, gde.New(app, fmt.Sprintf("INTERNAL ERROR: Could not create Procfile: %s", err.Error())))
			return err
		}

		em.Trigger(eventName, gde.New(app, "Procfile created!"))
		em.Trigger(eventName, gde.New(app, fmt.Sprintf("%s", config.Procfile)))
	} else {
		em.Trigger(eventName, gde.New(app, "WARNING: No Procfile configuration found. Deployment is most likely going to fail."))
	}
	return nil
}
