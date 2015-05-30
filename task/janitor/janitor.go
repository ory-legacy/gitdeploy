package janitor

import (
	"errors"
	"fmt"
	"github.com/ory-am/gitdeploy/task"
)

func Cleanup(dir string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Cleaning up...")
		if err := task.Exec(w, "", "rm", "-rf", dir); err != nil {
			return errors.New(fmt.Sprintf("Could not remove temp file: %s", err.Error()))
		}
		return nil
	}
}
