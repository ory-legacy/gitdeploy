package job

import (
	"fmt"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"github.com/go-errors/errors"
)

func GetLogs(app string) (string, error) {
	o, err := exec.Command("flynn", "-a", app, "logs").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(o)
}
