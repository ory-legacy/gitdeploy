package eco
import (
	"os/exec"
	"log"
	"runtime"
	"strings"
	"github.com/ory-am/gitdeploy/task/git"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task"
)

func IsGitAvailable() {
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("Git CLI is required but not installed or not in path.")
	}
}

func InitGit() {
	if err := git.Init(); err != nil {
		log.Fatalf("Could not init git: $s", err.Error())
	}
}

func InitFlynn(clusterConf string) {
	w := new(task.WorkerLog)
	if err := flynn.AddKey(w); err != nil {
		log.Fatalf("Could not init flynn: $s", err.Error())
	}
	log.Println("Adding flynn cluster...")
	log.Println(clusterConf)
	args := append([]string{"cluster", "add"}, strings.Split(clusterConf, " ")...)
	log.Printf("%s", args)
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
