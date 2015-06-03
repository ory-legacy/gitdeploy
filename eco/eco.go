package eco

import (
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"regexp"
	"net/url"
	"errors"
	"fmt"
)

func IsGitAvailable() {
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("Git CLI is required but not installed or not in path.")
	}
}

func InitFlynn(clusterConf string) {
	w := make(task.WorkerLog)
	if err := flynn.AddKey(w); err != nil {
		log.Fatalf("Could not init flynn: $s", err.Error())
	}
	log.Println("Adding flynn cluster...")
	args := append([]string{"cluster", "add"}, strings.Split(clusterConf, " ")...)
	if o, err := exec.Command("flynn", args...).CombinedOutput(); err != nil {
		log.Fatalf("Could not add cluster (status: %s) (output: %s) (args: %s)", err.Error(), o, args)
	} else {
		log.Printf("Adding cluster successful: %s", o)
	}
}

func IsFlynnAvailable() {
	_, err := exec.LookPath("flynn")
	if err != nil {
		if runtime.GOOS == "windows" {
			log.Fatal("Flynn CLI is required but not installed or not in path.")
		}
		log.Println("Could not find Flynn CLI, trying to install...")
		if o, err := exec.Command("sh", "bin/flynn-install.sh").CombinedOutput(); err != nil {
			log.Printf("Could not install Flynn CLI (%s): %s", err.Error(), o)
		} else if _, err := exec.LookPath("flynn"); err != nil {
			log.Fatal("Could not install Flynn CLI.")
		}
		log.Println("Flynn installed successfully!")
	}
}

func GetFlynnHost() (s string, err error) {
	reg := regexp.MustCompile(`(?mi)[a-z0-9\-A-Z]+\s+(https\:\/\/)controller\.([\.a-zA-Z0-9]+)\s+\(default\)$`)
	if o, err := exec.Command("flynn", "cluster").CombinedOutput(); err != nil {
		return "", err
	} else {
		s = string(o)
	}
	results := reg.FindStringSubmatch(s)
	if len(results) < 2 {
		return "", errors.New(fmt.Sprintf("Could not parse cluster information. Result: %s. Data: %s", results, s))
	} else {
		if u, err := url.Parse(results[1] + results[2]); err != nil {
			return "", err
		} else {
			return u.Host, nil
		}
	}
}
