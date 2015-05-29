package git

import "github.com/ory-am/gitdeploy/task"

// Run runs "git clone".
func Clone(w task.WorkerLog, repository, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Cloning repository...")
		if err := task.Exec(w, "git", "clone", "--progress", repository, wd); err != nil {
			return err
		}
		return nil
	}
}
