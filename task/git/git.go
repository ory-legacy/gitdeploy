package git

import (
	"fmt"
	"github.com/ory-am/gitdeploy/task"
	"os"
	"os/exec"
	"runtime"
)

type Git struct{ *task.Helper }

func (h *Git) CreateDirectory() (destination string) {
	destination = fmt.Sprintf("%s/%s", os.TempDir(), h.App)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s", os.TempDir(), h.App)
	}
	return destination
}

func (f *Git) Init() error {
	if err := exec.Command("git", "config", "user.name", "gitdeploy").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "config", "user.name", "gitdeploy").Run(); err != nil {
		return err
	}
	return nil
}

func (f *Git) Commit() error {
	if err := exec.Command("git", "commit", "-a", "-m", "gitdeploy").Run(); err != nil {
		return err
	}
	return nil
}

func (f *Git) AddAll() error {
	if err := exec.Command("git", "add", "--all").Run(); err != nil {
		return err
	}
	return nil
}
