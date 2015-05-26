package flynn

import (
	"errors"
	"fmt"
	"github.com/ory-am/gitdeploy/task"
	"net/url"
	"os/exec"
	"regexp"
)

type Flynn struct {
	envVars []string
	*task.Helper
}

func (h *Flynn) GetCluster() (*url.URL, error) {
	reg := regexp.MustCompile(`(?mi)[a-z0-9\-A-Z]+\s+(https\:\/\/)controller\.([\.a-zA-Z0-9]+)\s+\(default\)$`)
	out, err := exec.Command("flynn", "cluster").CombinedOutput()
	if err != nil {
		return nil, err
	}

	if results := reg.FindStringSubmatch(string(out)); len(results) < 2 {
		return nil, errors.New(fmt.Sprintf("Could not parse cluster information. Result: %s. Data: %s", results, out))
	} else {
		return url.Parse(results[1] + results[2])
	}
}

func (f *Flynn) AddEnvVar(key, value string) {
	if f.envVars == nil {
		f.envVars = make([]string, 0)
	}
	f.envVars = append(f.envVars, key+"="+value)
}

func (f *Flynn) CommitEnvVars() error {
	return exec.Command("flynn", "-a", f.App, "env", "set", f.envVars...).Run()
}

func (f *Flynn) GetLogs() (string, error) {
	o, err := exec.Command("flynn", "-a", f.App, "log").CombinedOutput()
	return string(o), err
}

func (f *Flynn) GetProcs() (string, error) {
	o, err := exec.Command("flynn", "-a", f.App, "ps").CombinedOutput()
	return string(o), err
}
