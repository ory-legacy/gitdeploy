package janitor

import (
	"github.com/ory-am/gitdeploy/task"
	"fmt"
	"github.com/ory-am/gitdeploy/Godeps/_workspace/src/github.com/go-errors/errors"
)

type Janitor struct {
	*task.Helper
}

func (j *Janitor) Cleanup() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(j.EventName, "Cleaning up...")
	if err := j.Exec(w, "../", "rm", "-rf", j.WorkingDirectory); err != nil {
		return w, errors.New(fmt.Sprintf("Could not remove temp file: %s", err.Error()))
	}
	return w, nil
}
