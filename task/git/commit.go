package git

import "github.com/ory-am/gitdeploy/task"

func Commit(w task.WorkerLog, wd string) func(task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		return task.Exec(w, wd, "git", "commit", "-a", "-m", "gitdeploy")
	}
}