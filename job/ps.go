package job

import (
	"os/exec"
)

func GetPS(app string) (string, error) {
	o, err := exec.Command("flynn", "-a", app, "ps").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(o), nil
}
