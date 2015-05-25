package git

import "github.com/ory-am/gitdeploy/task"

type CloneHelper struct {
	Repository string
	*task.Helper
}

// Run runs "git clone".
func (h *CloneHelper) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(h.EventName, "Cloning repository...")

	if err := h.Exec(w, "git", "clone", "--progress", h.Repository, h.WorkingDirectory); err != nil {
		return w, err
	}
	return w, nil
}
