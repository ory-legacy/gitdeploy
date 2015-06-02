package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CreateDirectory(app string) (destination string) {
	tempDir := strings.Trim(os.TempDir(), "/\\")
	destination = fmt.Sprintf("%s/%s", tempDir, app)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s", tempDir, app)
	}
	return
}

func Init() error {
	if err := exec.Command("git", "config", "user.name", "gd").Run(); err != nil {
		return err
	}
	return exec.Command("git", "config", "user.email", "gd@gitdeploy").Run()
}
