package flynn

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
)

type EnvHelper struct {
	envVars []string
}

func (f *EnvHelper) AddEnvVar(key, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	if f.envVars == nil {
		f.envVars = make([]string, 0)
	}
	f.envVars = append(f.envVars, key+"="+value)
}

func (f *EnvHelper) CommitEnvVars(app string) error {
	if len(f.envVars) == 0 {
		return nil
	}
	return exec.Command("flynn", append([]string{"-a", app, "env", "set"}, f.envVars...)...).Run()
}

func GetLogs(app string) (string, error) {
	o, err := exec.Command("flynn", "-a", app, "log").CombinedOutput()
	return string(o), err
}

func GetProcs(app string) (string, error) {
	o, err := exec.Command("flynn", "-a", app, "ps").CombinedOutput()
	return string(o), err
}

func GetCluster() (*url.URL, error) {
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
