package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CreateDirectory(app string) (destination string) {
	destination = fmt.Sprintf("%s/%s", os.TempDir(), app)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s", os.TempDir(), app)
	}
	return
}

func Init() error {
	if err := exec.Command("git", "config", "user.name", "gd").Run(); err != nil {
		return err
	}
	return exec.Command("git", "config", "user.email", "gd@gitdeploy").Run()
}
