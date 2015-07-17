package git

import "github.com/ory-am/gitdeploy/task"

// Checkout checks out a git repository
func Checkout(app, path, ref string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Checking out branch...")
		if err := task.Exec(w, path, "git", "checkout", "-b", app, ref); err != nil {
			return err
		}
		return nil
	}
}
