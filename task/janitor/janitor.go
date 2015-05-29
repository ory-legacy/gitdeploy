package janitor

import (
	"fmt"
	"errors"
	"github.com/ory-am/gitdeploy/task"
)

func Cleanup(w task.WorkerLog, dir string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Cleaning up...")
		if err := task.Exec(w, "", "rm", "-rf", dir); err != nil {
			return w, errors.New(fmt.Sprintf("Could not remove temp file: %s", err.Error()))
		}
		return w, nil
	}
}
